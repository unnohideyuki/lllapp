import React from 'react'
import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { createRoot } from 'react-dom/client'
import UserPage from './UserPage'
import ReadingRecord from './ReadingRecord'

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
        <Route path="/users/:id/" element={<UserPage />} />
      </Routes>
    </BrowserRouter>
  )
}
