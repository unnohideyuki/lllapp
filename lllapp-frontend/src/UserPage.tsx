import React, { useState, useEffect } from 'react'
import axios from 'axios'
import { BrowserRouter, Routes, Route, useParams } from 'react-router-dom'
import { SubmitHandler, useForm } from 'react-hook-form'

import './index.css'
const backendBaseURL = 'http://localhost:8080'

const UserPage = () => {
  const params = useParams()

  return (
    <div>
      <h1>{params.id}</h1>
    </div>
  )
}

export default UserPage
