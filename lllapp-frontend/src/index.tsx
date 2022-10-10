import React, { useState, useEffect } from 'react'
import axios from 'axios'
import { BrowserRouter, Routes, Route, useParams } from 'react-router-dom'
import { createRoot } from 'react-dom/client'
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
    for (let i = p0; i <= p1; i++)
      pages.push(i)
    return (
      <tr>
        <td className="section-name">{param.name}</td>
        <td>{pages.map((i) => renderSquare(progress[i]))}</td>
      </tr>
    )
  } else {
    return (
      <tr>
        <td colSpan={2} className="part-name">{param.name}</td>
      </tr>
    )
  }
}

const ReadingRecord = () => {
  const params = useParams()

  const [datas, setDatas] = useState<FetchData[]>([])
  const urlTOC = `${backendBaseURL}/users/${params.id}/books/${params.num}/toc`
  useEffect(() => { axios.get(urlTOC).then((res) => setDatas(res.data))}, [])

  const [progress, setDatas2] = useState<number[]>([])
  const urlProgress = `${backendBaseURL}/users/${params.id}/books/${params.num}/pghistory`
  useEffect(() => { axios.get(urlProgress).then((res) => setDatas2(res.data))}, [])

  return (
    <div>
      <div className="reading-record">
        <div className="reading-record-info">
          {/* todo */}
        </div>
      </div>
      <div className="reading-record-table">
        <table className="reading-record-table">
          <tbody>{datas.map((item: FetchData) => renderRow(item, progress))}</tbody>
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
        <Route path="/users/:id/books/:num/progress" element={<ReadingRecord />} />
      </Routes>
    </BrowserRouter>
  )
}
