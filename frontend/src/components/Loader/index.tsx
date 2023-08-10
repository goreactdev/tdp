import React from 'react'
import { BeatLoader } from 'react-spinners'

export const Loader = () => {
  return (
    <div className="flex h-[80vh] w-full items-center justify-center">
      <BeatLoader color="#0088CC" />{' '}
    </div>
  )
}
