import { useLayoutEffect } from 'react'
import {
  useNavigate,
  useParams,
} from 'react-router-dom'

import { Profile } from '../../components/Profile'
import SBTGrid from '../../components/SBTGrid/SBTGrid'
import { Container } from '../../components/SBTGrid/SBTGrid.styles'
import {
  useGetUserByUsernameQuery,
} from '../../services/users'

const ProfilePage = () => {
  const { username } = useParams<{ username?: string }>()
  const navigate = useNavigate()

  if (!username) {
    return <div>Username is undefined</div>
  }

  const { data, isSuccess, isLoading, isError } = useGetUserByUsernameQuery({
    username: username ?? '',
  })

  useLayoutEffect(() => {
    if (!isLoading && !data) {
      navigate('/404')
    }

    if (isError) {
      navigate('/404')
    }
  }, [data, isLoading, isError])

  return (
    <Container>
      {data && isSuccess && (
        <>
          <Profile user={data.user} />
          <SBTGrid />
        </>
      )}
    </Container>
  )
}

export default ProfilePage
