import React from 'react'
import { BsGithub, BsTelegram } from 'react-icons/bs'
import { Link, useNavigate, useSearchParams } from 'react-router-dom'

import Button from '../../components/Button/Button'
import Pagination from '../../components/Pagination/Pagination'
import { DefaultAvatar } from '../../components/Profile/Profile'
import { Table } from '../../components/Table'
import { PageHeaderText } from '../../components/Text'
import { useGetTopUsersQuery } from '../../services/users'

const RankingPage = () => {
  const [searchParams, setSearchParams] = useSearchParams()

  const pageSize = searchParams.get('pageSize') || '6'

  const [currentPage, setCurrentPage] = React.useState(
    parseInt(searchParams.get('current') || '1')
  )
  const navigate = useNavigate()

  const { data, isSuccess, isLoading } = useGetTopUsersQuery({
    end: currentPage * parseInt(pageSize),
    start: (currentPage - 1) * parseInt(pageSize),
  })

  React.useEffect(() => {
    if (!isLoading && !data) {
      navigate('/404')
    }
    setSearchParams({ current: currentPage.toString(), pageSize })
  }, [data, currentPage])

  return (
    <div className="flex w-full flex-col">
      <div className="flex justify-between">
        <PageHeaderText>TOP Developers</PageHeaderText>
        <Button
          href="https://ton-org.notion.site/TDP-Achievements-list-bc14d2b34ddb437d8019ac839cc03ea2"
          color="blue"
          className="scale-95"
        >
          How to Become #1
        </Button>
      </div>
      {data && isSuccess && (
        <Table
          isUserRating
          config={[
            {
              key: 'id',
              label: '#',
            },
            {
              key: 'user',
              label: 'User',
            },
            {
              key: 'rating',
              label: 'Rating',
            },
            {
              key: 'awards_count',
              label: 'Awards',
            },
            {
              key: 'linked_accounts',
              label: 'Links',
            },
          ]}
          data={data.map((user, i) => ({
            awards_count: user.awards_count || 0,
            id: currentPage * i + 1,
            linked_accounts: (
              <div className="flex space-x-2 text-gray-500">
                {user.linked_accounts?.filter(
                  (account) => account.provider === 'telegram'
                ).map((account) => {
                  
                    return (
                      <a href={`https://t.me/${account.login}`} target="_blank">
                        <BsTelegram className="h-5 w-5 transition-all duration-300 hover:scale-110" />
                      </a>
                    )
                  }
                )}
                {user.linked_accounts?.filter(
                  (account) => account.provider === 'github'
                ).map((account) => (
                  <a href={`https://github.com/${account.login}`} target="_blank">
                    <BsGithub className="h-5 w-5 transition-all duration-300 hover:scale-110" />
                  </a>
                ))}
              </div>
            ),
            rating: (
              <>
                {new Intl.NumberFormat('en-US', {
                  // compact 1k, 10k, 1m, 10m
                  notation: 'compact',
                }).format(user.rating || 0)}
              </>
            ),
            user: (
              <Link to={`/user/${user.username}`} className="flex items-center">
                {user.avatar_url ? (
                  <img
                    src={user.avatar_url}
                    alt={user.username}
                    className="h-12 w-12 rounded-full object-cover"
                  />
                ) : (
                  <DefaultAvatar size="small" />
                )}

                <div className="ml-4">
                  {user.first_name + ' ' + user.last_name}
                </div>
              </Link>
            ),
          }))}
        />
      )}

      {data && data?.length > 20 && (
        <Pagination
          currentPage={currentPage}
          setCurrentPage={setCurrentPage}
          pageSize={pageSize}
          count={0}
        />
      )}
    </div>
  )
}

export default RankingPage
