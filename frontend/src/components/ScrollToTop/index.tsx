import { useEffect } from 'react'
import { useLocation } from 'react-router-dom'

const ScrollToTop = () => {
  const { pathname } = useLocation()

  useEffect(() => {
    const container = document.getElementById('container')
    if (container) {
      container.scrollTop = 0
    }
  }, [pathname])

  return null
}

export default ScrollToTop
