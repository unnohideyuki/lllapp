import React from "react"
import { BrowserRouter, Routes, Route, useParams } from 'react-router-dom'
import { createRoot } from 'react-dom/client'
import ReadingRecord from "./ReadingRecord"

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
