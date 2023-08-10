import React from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'

import Button from '../../components/Button/Button'
import Pagination from '../../components/Pagination/Pagination'
import { Table } from '../../components/Table'
import { PageHeaderText } from '../../components/Text'
import {
  useGetAchievementsQuery,
  useUpdateAchievementMutation,
} from '../../services/users'

const AchievementsPage = () => {
  const [searchParams, setSearchParams] = useSearchParams()

  const pageSize = searchParams.get('pageSize') || '6'

  const [currentPage, setCurrentPage] = React.useState(
    parseInt(searchParams.get('current') || '1')
  )
  const navigate = useNavigate()

  const { data, isSuccess, isLoading } = useGetAchievementsQuery({
    end: currentPage * parseInt(pageSize),
    start: (currentPage - 1) * parseInt(pageSize),
  })

  const [updateApproved] = useUpdateAchievementMutation()

  React.useEffect(() => {
    if (!isLoading && !data) {
      navigate('/404')
    }
    setSearchParams({ current: currentPage.toString(), pageSize })
  }, [data, currentPage])

  return (
    <div className="flex w-full flex-col">
      <PageHeaderText>Incoming Achievements</PageHeaderText>
      {data && isSuccess && !data?.achievements && (
        <div className="flex w-full flex-col items-center justify-center">
          <p className="text-2xl font-bold">No achievements</p>
        </div>
      )}

      {data && isSuccess && data.achievements && (
        <Table
          config={[
            {
              key: 'id',
              label: '#',
            },
            {
              key: 'name',
              label: 'Name',
            },

            {
              key: 'description',
              label: 'Description',
            },
            {
              key: 'image',
              label: 'Image',
            },
            {
              key: 'weight',
              label: 'Weight',
            },
            {
              key: 'status',
              label: 'Status',
            },
          ]}
          data={data.achievements.map((sbt) => ({
            description: sbt.description,
            id: sbt.achievement_id,
            image: (
              <div className="h-20 w-20">
                <img
                  src={sbt.image_url}
                  alt={sbt.name}
                  className="h-20 w-20 rounded-2xl object-cover"
                />
              </div>
            ),
            name: sbt.name,
            status: (
              <>
                {!sbt.processed && !sbt.approved && (
                  <Button
                    onClick={() => {
                      updateApproved({
                        approved_by_user: true,
                        id: sbt.achievement_id,
                      })
                    }}
                    color="blue"
                  >
                    Approve
                  </Button>
                )}
                {sbt.processed && sbt.approved && (
                  <Button disabled color="white">
                    Minted
                  </Button>
                )}
                {!sbt.processed && sbt.approved && (
                  <Button disabled color="yellow">
                    In Progress
                  </Button>
                )}
              </>
            ),
            weight: (
              <>
                {new Intl.NumberFormat('en-US', {
                  compactDisplay: 'short',
                }).format(sbt.weight || 0)}
              </>
            ),
          }))}
        />
      )}

      {data && data.count > 20 && (
        <Pagination
          currentPage={currentPage}
          setCurrentPage={setCurrentPage}
          pageSize={pageSize}
          count={data.count}
        />
      )}
    </div>
  )
}

export default AchievementsPage
