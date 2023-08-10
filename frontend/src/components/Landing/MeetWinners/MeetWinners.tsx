import { BsGithub, BsTelegram } from 'react-icons/bs'

import { useGetTopUsersQuery } from '../../../services/users'
import { SocialSvg, SocialLink } from '../../Icons/Icons.styles'
import { DefaultAvatar } from '../../Profile/Profile'

import {
  ExploreLink,
  Name,
  Position,
  ProfilePicture,
  ProfilePictureWrapper,
  SocialLinksWrapper,
  Title,
  WinnerCard,
  WinnersGrid,
} from './MeetWinners.styles'

export const jobs = [
  'Funemployed',
  'Professional Netflix Researcher',
  'Domestic CEO',
  'Between Successes',
  'Master of Leisure Arts',
  'Job-free by Choice',
  'Adventurer of Unchartered Territories',
  'Undiscovered Talent Magnet',
  'Full-Time Dream Chaser',
  'Certified Couch Explorer',
]

const MeetWinners = () => {
  const { data, isSuccess } = useGetTopUsersQuery({
    end: 8,
    start: 0,
  })

  return (
    <div>
      <div className="mb-10 mt-16 flex flex-col items-center justify-between space-y-3 lg:flex-row lg:space-y-0">
        <Title>Most Active Members</Title>
        <ExploreLink to="/rewards">
          Explore all members
          <span>
            <SocialSvg
              width="26"
              height="30"
              viewBox="0 0 24 25"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                d="M8 18L13 13.5L8 9"
                stroke="currentColor"
                strokeWidth="2.2"
                strokeLinecap="round"
                strokeLinejoin="round"
              />
            </SocialSvg>
          </span>
        </ExploreLink>
      </div>

      <WinnersGrid>
        {isSuccess &&
          data &&
          data.map((winner) => (
            <WinnerCard to={`user/${winner.username}`} key={winner.username}>
              {winner.avatar_url ? (
                <ProfilePictureWrapper>
                  <ProfilePicture
                    src={winner.avatar_url}
                    loading="lazy"
                    alt={winner.first_name + ' ' + winner.last_name}
                  />
                </ProfilePictureWrapper>
              ) : (
                <div className="mb-4">
                  <DefaultAvatar size="medium" />
                </div>
              )}

              <div>
                <Name>{winner.first_name + ' ' + winner.last_name}</Name>
                <Position>
                  {winner.job ||
                    jobs[
                      winner.id > 9 ? Number(String(winner.id)[0]) : winner.id
                    ]}
                </Position>

                <SocialLinksWrapper>
                  <div className="flex gap-4">
                    {winner.linked_accounts?.filter(
                      (account) => account.provider === 'telegram'
                    ).map((account) => (
                      <SocialLink
                        key={account.login}
                        onClick={(e) => e.stopPropagation()}
                        href={
                            `https://t.me/${account.login}`
                        }
                        target="_blank"
                      >
                          <BsTelegram className="h-5 w-5" />
                      </SocialLink>
                    ))}
                    {winner.linked_accounts?.filter(
                      (account) => account.provider === 'github'
                    ).map((account) =>
                      <SocialLink
                        onClick={(e) => e.stopPropagation()}
                        href={`https://github.com/${account.login}`}
                        target="_blank"
                      >
                        <BsGithub className="h-5 w-5" />
                      </SocialLink>
                    )}
                  </div>
                </SocialLinksWrapper>
              </div>
            </WinnerCard>
          ))}
      </WinnersGrid>
    </div>
  )
}

export default MeetWinners
