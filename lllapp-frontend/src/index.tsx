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
  let bg = '#fff'

  if (i == undefined) i = 0
  if (i > 2) bg = '#fc7f00'
  else if (i == 2) bg = '#ffbd00'
  else if (i == 1) bg = '#e4ff7a'

  const style = { background: bg }

  return <div className="square" style={style}></div>
}

const renderRow = (param: FetchData, progress: number[]): JSX.Element => {
  const pages = []
  if (Number(param.page_from) >= 0) {
    for (let i = Number(param.page_from); i <= Number(param.page_to); i++)
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
        <td colSpan={2} className="part-name">
          {param.name}
        </td>
      </tr>
    )
  }
}

const ReadingRecord = () => {
  const [datas, setDatas] = useState<FetchData[]>([])
  const params = useParams()
  const urlAPI = `${backendBaseURL}/users/${params.id}/books/${params.num}/toc`
  useEffect(() => { axios.get(urlAPI).then((res) => setDatas(res.data))}, [])

  return (
    <div>
      <div className="reading-record">
        <div className="reading-record-info">
          {/* todo */}
        </div>
      </div>
      <div className="reading-record-table">
        <table className="reading-record-table">
          <tbody>{datas.map((item: FetchData) => renderRow(item, []))}</tbody>
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
