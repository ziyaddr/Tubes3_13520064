import React from 'react'
import { Link } from 'react-router-dom'

const Landing = () => {
  const reload = (e) => {
    e.preventDefault();
    window.location.reload(true)
  }
  return (
    <div>
      <p>Welcome to DNA Application</p>
      <div onClick={reload}>
        <p><Link to="/add-disease">Add Disease</Link></p>
        <p><Link to="/disease-prediction">Disease Prediction</Link></p>
        <p><Link to="/detail-prediction">Detail Prediction</Link></p>
      </div>
    </div>
  )
}

export default Landing
