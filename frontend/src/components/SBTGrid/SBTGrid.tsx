import { useEffect, useState } from 'react'
import { useNavigate, useParams, useSearchParams } from 'react-router-dom'

import {
  Badge,
  Card,
  CardContent,
  CardSubtitle,
  CardTitle,
  Content,
  Grid,
  ImageContent,
  Img,
  InnerContainer,
  TextContainer,
} from '../../components/SBTGrid/SBTGrid.styles'
import { useMemoizedUser } from '../../hooks/useMemoizedUser'
import {
  usePinNftMutation,
  useGetNftsByUsernameQuery,
} from '../../services/users'
import Pagination from '../Pagination/Pagination'
import { PageHeaderText } from '../Text'

const SBTGrid = () => {
  const { username } = useParams<{ username?: string }>()

  const { user } = useMemoizedUser()

  const [triggerPin] = usePinNftMutation()

  const [searchParams, setSearchParams] = useSearchParams()

  const pageSize = searchParams.get('pageSize') || '6'

  const [currentPage, setCurrentPage] = useState(
    parseInt(searchParams.get('current') || '1')
  )
  const navigate = useNavigate()

  const {
    data,
    isSuccess: isSuccessNfts,
    isLoading: isLoadingNfts,
  } = useGetNftsByUsernameQuery({
    end: currentPage * parseInt(pageSize),
    start: (currentPage - 1) * parseInt(pageSize),
    username: username ?? '',
  })

  useEffect(() => {
    if (!isLoadingNfts && !data) {
      navigate('/404')
    }
    setSearchParams({ current: currentPage.toString(), pageSize })
  }, [data, currentPage])

  return (
    <Content>
      <InnerContainer>
        <PageHeaderText className="mb-4">SBT tokens</PageHeaderText>

        <Grid>
          {isLoadingNfts && <div>Loading...</div>}
          {data?.nfts.length === 0 && (
            <Card>
              <ImageContent>
                <Img
                  src={
                    'https://images.unsplash.com/photo-1609743522653-52354461eb27?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=987&q=80'
                  }
                  loading="lazy"
                  alt="Photo by Austin Wade"
                />
              </ImageContent>
              <CardContent>There is no SBT token yet.</CardContent>
            </Card>
          )}
          {data?.nfts &&
            data.nfts.map((nft) => (
              <Card>
                <ImageContent>
                  <a
                    target="_blank"
                    href={`https://getgems.io/nft/${nft.friendly_address}`}
                  >
                    <Img
                      src={nft.image}
                      loading="lazy"
                      alt="Photo by Austin Wade"
                    />
                  </a>
                  {!nft.is_pinned &&
                    user?.friendly_address === nft.friendly_owner_address && (
                      <Badge
                        onClick={(e) => {
                          e.stopPropagation()

                          triggerPin({ id: nft.id })
                        }}
                        color="blue"
                      >
                        PIN
                      </Badge>
                    )}
                  {nft.is_pinned &&
                    user?.friendly_address === nft.friendly_owner_address && (
                      <Badge
                        onClick={(e) => {
                          e.stopPropagation()
                          triggerPin({ id: nft.id })
                        }}
                        color="red"
                      >
                        PINNED
                      </Badge>
                    )}
                </ImageContent>

                <CardContent
                  target="_blank"
                  href={`https://getgems.io/nft/${nft.friendly_address}`}
                >
                  <TextContainer>
                    <CardTitle>
                      {nft.name.length > 20
                        ? `${nft.name.slice(0, 20)}...`
                        : nft.name}
                    </CardTitle>
                    <CardSubtitle>
                      {nft.description.length > 50
                        ? `${nft.description.slice(0, 50)}...`
                        : nft.description}
                    </CardSubtitle>
                  </TextContainer>
                  <span className="text-sm text-green-500">+{nft.weight}</span>
                </CardContent>
              </Card>
            ))}
        </Grid>
      </InnerContainer>
      {data && isSuccessNfts && data.count > 9 && (
        <Pagination
          currentPage={currentPage}
          setCurrentPage={setCurrentPage}
          pageSize={pageSize}
          count={data.count}
        />
      )}
    </Content>
  )
}

export default SBTGrid
