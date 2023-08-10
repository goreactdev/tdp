import React from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'

import Hero from '../../components/Landing/Hero'
import MeetWinners from '../../components/Landing/MeetWinners'
import Stats from '../../components/Landing/Stats'
import { useGetTopUsersQuery } from '../../services/users'

const LandingPage = () => {
  return (
    <>
      <div>
        <Hero />
        <Stats />
      </div>

      <MeetWinners />
    </>
  )
}

export default LandingPage
