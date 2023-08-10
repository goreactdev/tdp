import { AiFillGithub, AiFillTrophy } from 'react-icons/ai'
import { BsFillPersonFill, BsTelegram } from 'react-icons/bs'
import { TbBrandTelegram } from 'react-icons/tb'

import type { User } from '../../services/types'
import Button from '../Button'

import { jobs } from '../Landing/MeetWinners/MeetWinners'

import {
  Avatar,
  AvatarContainer,
  AwardCount,
  AwardCountText,
  AwardsContainer,
  AwardsText,
  AwardsTitle,
  CardContainer,
  CardFooter,
  CardHeader,
  CertificationsBlock,
  CertificationsContainer,
  CertificationsTitle,
  DescriptionContainer,
  DescriptionText,
  DescriptionTitle,
  Diamond,
  DiamondContainer,
  EmptyDiamond,
  FilledDiamond,
  InfoContainer,
  InfoRow,
  JobTitle,
  LanguagesBlock,
  LanguagesContainer,
  LanguagesTitle,
  LinkedAccount,
  LinkedAccountIcon,
  LinkedAccountText,
  LinkedAccountsContainer,
  LinkedAccountsTitle,
  Name,
  UserIcon,
  Username,
  Wrapper,
} from './Profile.styles'
import { BASE_URL } from '../../utils/config'
import React, { useEffect } from 'react'
import Popup from 'reactjs-popup'

const DiamondRating = ({ awards }: { awards: number }) => {
  const MAX_DIAMONDS = 5

  const filledDiamonds = Math.min(Math.ceil(awards / 500), MAX_DIAMONDS)
  const emptyDiamonds = MAX_DIAMONDS - filledDiamonds

  return (
    <DiamondContainer>
      {[...Array(filledDiamonds)].map((_, index) => (
        <Diamond key={index} aria-hidden="true">
          <FilledDiamond />
        </Diamond>
      ))}
      {[...Array(emptyDiamonds)].map((_, index) => (
        <Diamond key={index + filledDiamonds} aria-hidden="true">
          <EmptyDiamond />
        </Diamond>
      ))}
      <AwardCount>
        {filledDiamonds}{' '}
        <AwardCountText>
          (
          {new Intl.NumberFormat('en-US', { notation: 'compact' }).format(
            awards
          )}{' '}
          points)
        </AwardCountText>
      </AwardCount>
    </DiamondContainer>
  )
}

type DefaultAvatarProps = {
  size: 'small' | 'medium' | 'large'
}

export const DefaultAvatar: React.FC<DefaultAvatarProps> = ({ size }) => (
  <AvatarContainer>
    <UserIcon size={size} />
  </AvatarContainer>
)

export const Profile = ({ user }: { user: User }) => {
  const [open, setOpen] = React.useState(false)
  const closeModal = () => setOpen(false)


  return (
    <Wrapper>
      <CardContainer>
        <Popup
        open={open}
        onOpen={() => {
          setOpen(true)
          navigator.clipboard.writeText(`${BASE_URL}/user/${user.username}`)
          setTimeout(() => {
            setOpen(false)
          }, 1000)
        }}

        onClose={closeModal}
          trigger={
            <CardHeader>
              {user.avatar_url ? (
                <Avatar
                  className="mt-6 "
                  size="large"
                  imageUrl={user.avatar_url ?? ''}
                />
              ) : (
                <div className="mt-6">
                  <DefaultAvatar size="large" />
                </div>
              )}

              <div className="flex flex-col  items-center">
                <Name>{`${user.first_name} ${user.last_name}`}</Name>
                <Username>{`@${user.username}`}</Username>
              </div>
              <JobTitle>
                {user.job
                  ? user.job
                  : jobs[user.id > 9 ? Number(String(user.id)[0]) : user.id]}
              </JobTitle>
              <DiamondRating awards={user.rating ?? 0} />
              {user.linked_accounts &&
                user.linked_accounts.some(
                  (linked_account) => linked_account.provider === 'telegram'
                ) && (
                  <Button
                    href={
                      'https://t.me/' +
                      user.linked_accounts.find(
                        (linked_account) =>
                          linked_account.provider === 'telegram'
                      )?.login
                    }
                    color="blue"
                  >
                    Contact me
                  </Button>
                )}
            </CardHeader>
          }
          position="center center"
        >
          <div className='bg-green-500 rounded-2xl p-4 text-sm text-white font-bold bg-opacity-80'>Link is copied!</div>
        </Popup>

        <CardFooter>
          <InfoContainer>
            <InfoRow>
              <div className="flex items-center">
                <BsFillPersonFill className="fill-backgroundBlack" />
                <span className="ml-2">Member since</span>
              </div>
              <span className="font-semibold">
                {new Intl.DateTimeFormat('gb-GB', {
                  month: 'numeric',
                  year: 'numeric',
                }).format(new Date(user.created_at * 1000))}
              </span>
            </InfoRow>
            <InfoRow>
              <div className="flex items-center">
                <AiFillTrophy className="fill-backgroundBlack" />
                <span className="ml-2">Nr. of awards</span>
              </div>
              <span className="font-semibold">{user.awards_count ?? '0'}</span>
            </InfoRow>
            <InfoRow>
              <div className="flex items-center">
                <TbBrandTelegram className="fill-backgroundBlack" />
                <span className="ml-2">Last awards</span>
              </div>
              <span className="font-semibold">
                {user.last_award_at
                  ? new Intl.DateTimeFormat('gb-GB', {
                      month: 'numeric',
                      year: 'numeric',
                    }).format(new Date(user.last_award_at * 1000))
                  : 'Never'}
              </span>
            </InfoRow>
          </InfoContainer>
        </CardFooter>
      </CardContainer>

      <CardContainer>
        <DescriptionContainer>
          <DescriptionTitle>Description</DescriptionTitle>
          <DescriptionText>{user.bio || 'No description'}</DescriptionText>
        </DescriptionContainer>

        <LanguagesContainer>
          <LanguagesTitle>Languages</LanguagesTitle>
          <LanguagesBlock>
            {user.languages &&
              user.languages.length > 0 &&
              user.languages?.map((language) => <div>{language}</div>)}
            {user.languages && user.languages.length === 0 && 'No languages'}
          </LanguagesBlock>
        </LanguagesContainer>

        <LinkedAccountsContainer>
          <LinkedAccountsTitle>Linked Accounts</LinkedAccountsTitle>

          {user.linked_accounts &&
            user.linked_accounts.length > 0 &&
            user.linked_accounts?.map((linked_account) => (
              <LinkedAccount key={linked_account.provider}>
                <LinkedAccountIcon>
                  {linked_account.provider === 'github' && <AiFillGithub />}
                  {linked_account.provider === 'telegram' && <BsTelegram />}
                </LinkedAccountIcon>
                <LinkedAccountText
                  href={
                    linked_account.provider === 'github'
                      ? `https://github.com/${linked_account.login}`
                      : `https://t.me/${linked_account.login}`
                  }
                  target="_blank"
                >
                  {linked_account.provider[0].toUpperCase() +
                    linked_account.provider.slice(1)}
                </LinkedAccountText>
              </LinkedAccount>
            ))}

          {!user.linked_accounts && (
            <div className="text-base">No accounts</div>
          )}
        </LinkedAccountsContainer>

        <CertificationsContainer>
          <CertificationsTitle>Certifications</CertificationsTitle>
          <CertificationsBlock>
            {user.certifications &&
              user.certifications.length > 0 &&
              user.certifications?.map((certification) => certification)}
            {user.certifications &&
              user.certifications.length === 0 &&
              'No certifications'}
          </CertificationsBlock>
        </CertificationsContainer>
      </CardContainer>
    </Wrapper>
  )
}
