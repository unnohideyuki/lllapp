import React, { useState, useEffect } from 'react'
import axios from 'axios'
import { BrowserRouter, Routes, Route, useParams } from 'react-router-dom'
import { createRoot } from 'react-dom/client'
import { SubmitHandler, useForm } from 'react-hook-form'
import './index.css'

const backendBaseURL = 'http://localhost:8080'

type FetchData = {
  name: string
  page_from: string
  page_to: string
}

const renderSquare = (i: number | undefined): JSX.Element => {
  if (i == undefined) i = 0
  if (i > 4) i = 4
  const colors = ['#eeeeee', '#d6e685', '#8cc665', '#44a340', '#1e6823']
  const style = { background: colors[i] }

  return <div className="square" style={style}></div>
}

const renderRow = (param: FetchData, progress: number[]): JSX.Element => {
  const pages = []
  const p0 = Number(param.page_from)
  const p1 = Math.max(p0, Number(param.page_to))
  if (p0 >= 0) {
    for (let i = p0; i <= p1; i++) pages.push(i)
    return (
      <tr>
        <td className="section-name">{param.name}</td>
        <td>{pages.map((i) => renderSquare(progress[i]))}</td>
      </tr>
    )
  } else {
    return (
      <tr>
        <td colSpan={2} className="part-name">
          {param.name}
        </td>
      </tr>
    )
  }
}

type FormData = {
  pgfrom: string
  pgto: string
}

type BookInfo = {
  title: string
  author: string
  isbn: string
}

interface PanelState {
  mode: number
}

const ReadingRecord = () => {
  const [panelState, setPanelState] = useState<PanelState>({
    mode: 0,
  })

  const params = useParams()

  // Fetching TOC API
  const [datas, setDatas] = useState<FetchData[]>([])
  const urlTOC = `${backendBaseURL}/users/${params.id}/books/${params.num}/toc`
  useEffect(() => {
    axios.get(urlTOC).then((res) => setDatas(res.data))
  }, [])

  // Fetching Progress API
  const [progress, setDatas2] = useState<number[]>([])
  const urlProgress = `${backendBaseURL}/users/${params.id}/books/${params.num}/pghistory`
  useEffect(() => {
    axios.get(urlProgress).then((res) => setDatas2(res.data))
  }, [])

  // Fetching Info API
  const [info, setDatas3] = useState<BookInfo>({
    title: '',
    author: '',
    isbn: '',
  })
  const urlInfo = `${backendBaseURL}/users/${params.id}/books/${params.num}/info`
  useEffect(() => {
    axios.get(urlInfo).then((res) => setDatas3(res.data))
  }, [])

  const urlPostPages = `${backendBaseURL}/users/${params.id}/books/${params.num}/postpages`
  const { handleSubmit, register, reset } = useForm<FormData>()

  const postPages: SubmitHandler<FormData> = (data) => {
    console.log(data)
    const params = new URLSearchParams(data)
    axios
      .post(urlPostPages, params)
      .then((res) => {
        console.log(res)
      })
      .catch((e) => {
        console.error(e)
      })
    reset()
    window.location.reload()
  }

  // Fetching Last-Read pages
  const [lastpg, setDatas4] = useState<string[]>([])
  const urlLastRead = `${backendBaseURL}/users/${params.id}/books/${params.num}/lastread`
  useEffect(() => {
    axios.get(urlLastRead).then((res) => setDatas4(res.data))
  }, [])

  const changeMode = (x: number) => {
    setPanelState({ ...panelState, mode: x })
  }

  const renderPanel = () => {
    if (panelState.mode == 1) {
      return (
        <form onSubmit={handleSubmit(postPages)}>
          <div>
            <button type="button" onClick={() => changeMode(0)}>
              ✖️
            </button>
            <span>読んだページを入力：</span>
            <input
              type="number"
              {...register('pgfrom')}
              placeholder="page from"
            />
            <span>-</span>
            <input type="number" {...register('pgto')} placeholder="page to" />
            <span></span>
            <button type="submit">Send</button>
          </div>
        </form>
      )
    } else {
      // TODO: まだ１ページも読んでいない場合の対処
      return (
        <form>
          <button type="button" onClick={() => changeMode(1)}>
            ＋
          </button>
          <span>
            最後に読んだページ： page {lastpg[0]} - {lastpg[1]}（{lastpg[2]}）
          </span>
        </form>
      )
    }
  }

  return (
    <div>
      <div className="reading-record-info">
        <ul className="book-info">
          <li className="book-title">{info.title}</li>
          <li>{info.author}</li>
        </ul>
      </div>
      <div className="reading-record-panel">{renderPanel()}</div>
      <div className="reading-record-table">
        <table>
          <tbody>
            {datas.map((item: FetchData) => renderRow(item, progress))}
          </tbody>
        </table>
      </div>
    </div>
  )
}

// ========================================

const container = document.getElementById('root')
if (container) {
  const root = createRoot(container)
  root.render(
    <BrowserRouter>
      <Routes>
        <Route
          path="/users/:id/books/:num/progress"
          element={<ReadingRecord />}
        />
      </Routes>
    </BrowserRouter>
  )
}
