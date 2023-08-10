import { Popover } from '@headlessui/react'
import {
  TonConnectButton,
  useTonConnectUI,
} from '@tonconnect/ui-react'
import {  useState, useEffect } from 'react'
import {  AiOutlineClose } from 'react-icons/ai'
import { FaUserAlt } from 'react-icons/fa'
import { GiAchievement } from 'react-icons/gi'
import { GrFormDown } from 'react-icons/gr'
import { useDispatch } from 'react-redux'
import { Link } from 'react-router-dom'

import { logout } from '../../features/auth/authSlice'
import { useMemoizedUser } from '../../hooks/useMemoizedUser'
import {
  useCheckProofInBackendMutation,
  useGetProofPayloadQuery,
} from '../../services/users'
import Button from '../Button'
import type { PopoverItem } from '../Popover'
import PopoverContainer from '../Popover'
import PopoverElement from '../Popover'
import { Avatar } from '../Profile/Profile.styles'

import {
  ButtonContainer,
  HeaderContainer,
  LogoIcon,
  LogoLink,
  MobileButton,
  MobileIcon,
  NavContainer,
} from './Header.styles'

type NavButtonProps =
  | {
      items: PopoverItem[]
      href?: never
      to?: never
    }
  | {
      items?: never
      href: string
      to?: never
    }
  | {
      items?: never
      href?: never
      to: string
    }

type NavItem = {
  name: string
} & NavButtonProps

// array of nav items and buttons props

const navList = [
  // {
  //   name: 'Documentation',
  //   items: [
  //     {
  //       name: 'Getting Started',
  //       description: 'Explore more about TDP',
  //       href: 'https://docs.ton.org/',
  //       type: 'external_link',
  //       icon: (
  //         <div className="rounded-2xl bg-backgroundGray p-3">
  //           <ImBook
  //             className="
  //       h-6
  //       w-6
  //       fill-mainColor

  //       "
  //           />
  //         </div>
  //       ),
  //     },
  //     {
  //       name: 'Benefits of TDP',
  //       description: 'Learn how to use TDP',
  //       href: 'https://docs.ton.org/',
  //       type: 'external_link',
  //       icon: (
  //         <div className="rounded-2xl bg-backgroundGray p-3">
  //           {' '}
  //           <FaAward
  //             className="
  //       h-6
  //       w-6
  //       fill-mainColor

  //       "
  //           />
  //         </div>
  //       ),
  //     },
  //     {
  //       name: 'Ranking System',
  //       description: 'What does the rank mean',
  //       href: 'https://docs.ton.org/',
  //       type: 'external_link',
  //       icon: (
  //         <div className="rounded-2xl bg-backgroundGray p-3">
  //           {' '}
  //           <GiRank1
  //             className="
  //       h-6
  //       w-6
  //       fill-mainColor

  //       "
  //           />
  //         </div>
  //       ),
  //     },
  //     {
  //       name: 'Telegram BOT',
  //       description: 'How to use the Telegram BOT',
  //       href: 'https://docs.ton.org/',
  //       type: 'external_link',
  //       icon: (
  //         <div className="rounded-2xl bg-backgroundGray p-3">
  //           {' '}
  //           <BsTelegram
  //             className="
  //       h-6
  //       w-6
  //       fill-mainColor

  //       "
  //           />
  //         </div>
  //       ),
  //     },
  //   ],
  // },
  {
    name: 'Participate',
    href: 'https://ton-org.notion.site/TDP-Achievements-list-bc14d2b34ddb437d8019ac839cc03ea2',
  },

  {
    name: 'Leaderboard',
    to: '/rewards',
  },
  {
    href: 'https://ton-org.notion.site/How-to-get-rewarded-ad8ab607478d4a7ab8658051d4ce5bf7',
    name: 'Rewards',
  },
  {
    href: 'https://ton-org.notion.site/How-to-join-the-TDP-037f038bdb2848b8b46743bb38c7a473',
    name: 'Join as a partner',
  },
] as NavItem[]

const Header = () => {
  const dispatch = useDispatch()

  const { user } = useMemoizedUser()
  const [tonConnectUI] = useTonConnectUI()

  const profileItems = [
    {
      description: '',
      href: `/user/${user?.username}`,
      icon: null,
      name: 'My Profile',

      type: 'link',
    },
    {
      description: '',
      href: '/settings',
      icon: null,
      name: 'Settings',

      type: 'link',
    },

    {
      description: '',
      icon: null,
      name: 'Log Out',
      onClick: () => {
        dispatch(logout())
        tonConnectUI.disconnect()
        location.reload()
      },
      type: 'button',
    },
  ] as PopoverItem[]

  const [isMenuOpen, setIsMenuOpen] = useState(false)

  const { data: tonProofPayloadPromise } = useGetProofPayloadQuery()

  const [triggerCheckProof, { isLoading }] = useCheckProofInBackendMutation()

  const handleMenuToggle = () => {
    setIsMenuOpen(!isMenuOpen)
    // set to window owerflow hidden
  }

  useEffect(() => {
    if (isMenuOpen) {
      // get by id and set overflow hidden
      document
        .getElementById('container')
        ?.style.setProperty('overflow', 'hidden')
    }
    if (!isMenuOpen) {
      document
        .getElementById('container')
        ?.style.setProperty('overflow', 'auto')
    }

    return () => {
      document
        .getElementById('container')
        ?.style.setProperty('overflow', 'auto')
    }
  }, [isMenuOpen])

  if (!tonProofPayloadPromise) {
    tonConnectUI.setConnectRequestParameters(null)
  } else {
    tonConnectUI.setConnectRequestParameters({
      state: 'ready',
      value: { tonProof: tonProofPayloadPromise.payload },
    })
  }

  useEffect(() => {
    tonConnectUI.onStatusChange(async (wallet) => {
      if (
        wallet?.connectItems?.tonProof &&
        'proof' in wallet.connectItems.tonProof
      ) {
        triggerCheckProof({
          account: wallet.account,
          profItemReply: wallet.connectItems.tonProof,
        })
      }
    })
  }, [])

  const isMobile = window.innerWidth < 768

  // tonConnectUI.disconnect()
  return (
    <HeaderContainer>
      <div className="mx-auto flex w-full max-w-screen-xl justify-between px-4 lg:px-8">
        <LogoLink to="/" aria-label="logo">
          <LogoIcon
            width="95"
            height="94"
            viewBox="0 0 95 94"
            fill="currentColor"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path d="M96 0V47L48 94H0V47L48 0H96Z" />
          </LogoIcon>
          TDP
        </LogoLink>

        <NavContainer>
          {navList.map((navItem) => (
            <NavButton
              href={navItem.href}
              to={navItem.to}
              items={navItem.items}
              key={navItem.name}
            >
              {navItem.name}
            </NavButton>
          ))}
        </NavContainer>

        {!user && !isMobile && (
          <ButtonContainer>
            {!isLoading && <TonConnectButton className="h-[45px] w-[167px]" />}
          </ButtonContainer>
        )}

        {user && (
          <div className="hidden items-center justify-between lg:flex">
            <Link
              to={'/achievements'}
              className="mr-5 flex cursor-pointer space-x-4 rounded-full bg-mainColor p-[6px] transition-all duration-300 hover:brightness-125"
            >
              <GiAchievement size={16} className="h-6 w-6 fill-white" />
            </Link>

            <PopoverElement
              header={
                <div className="border-b bg-white px-4 py-2">
                  <Link to={`/user/${user.username}`}>
                    <div className=" text-sm font-medium">
                      {user.first_name + ' ' + user.last_name}
                    </div>
                    <div className="text-sm">
                      {user.friendly_address.slice(0, 6) +
                        '...' +
                        user.friendly_address.slice(-6)}
                    </div>
                  </Link>
                </div>
              }
              items={profileItems}
            >
              <Popover.Button className="flex cursor-pointer items-center font-medium outline-none">
                {user.avatar_url ? (
                  <Avatar size="small" imageUrl={user.avatar_url} />
                ) : (
                  <div className="rounded-full bg-mainColor transition-transform duration-300 hover:scale-105">
                    <FaUserAlt size={35} className="p-2" fill="white" />
                  </div>
                )}
                <GrFormDown className="ml-1" size={20} />
              </Popover.Button>
            </PopoverElement>
          </div>
        )}

        <div
          className={`fixed left-0 top-0 z-10 mt-16  flex h-screen w-full flex-col bg-white transition-all duration-300 ${
            !isMenuOpen ? 'hidden ' : 'block'
          }`}
        >
          {user && (
            <>
              <Link
                onClick={() => {
                  setIsMenuOpen(false)
                }}
                className="p-6 font-medium"
                to={`/user/${user.username}`}
              >
                My Profile
              </Link>

              <Link
                onClick={() => {
                  setIsMenuOpen(false)
                }}
                className="p-6 font-medium"
                to="/settings"
              >
                Settings
              </Link>

              <Link
                onClick={() => {
                  setIsMenuOpen(false)
                }}
                className="p-6 font-medium"
                to="/achievements"
              >
                Achievements
              </Link>
            </>
          )}
          <a
            href="https://docs.ton.org"
            target="_blank"
            className="p-6 font-medium"
          >
            Documentation
          </a>
          <Link
            to="/rewards"
            onClick={() => {
              setIsMenuOpen(false)
            }}
            className="p-6 font-medium"
          >
            Rewards
          </Link>
          <a
            href="https://ton.org"
            className="border-b border-gray-300 p-6 font-medium"
          >
            TON
          </a>
          {isMobile && !user && (
            <div className="mt-4 flex justify-center">
              <TonConnectButton className="h-[45px] w-[167px]" />
            </div>
          )}

          {user && (
            <div className="mt-4  flex justify-center">
              <Button
                onClick={() => {
                  dispatch(logout())
                  tonConnectUI.disconnect()
                  location.reload()
                  setIsMenuOpen(false)
                }}
                className="scale-125"
                color="blue"
              >
                Logout
              </Button>
            </div>
          )}
        </div>

        <MobileButton onClick={handleMenuToggle}>
          <MobileIcon
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 20 20"
            fill="currentColor"
            className={
              'absolute right-5 transition-all duration-300 ' +
              (isMenuOpen ? 'opacity-0' : 'opacity-100')
            }
          >
            <path
              fillRule="evenodd"
              d="M3 5a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zM3 10a1 1 0 011-1h6a1 1 0 110 2H4a1 1 0 01-1-1zM3 15a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z"
              clipRule="evenodd"
            />
            Menu
          </MobileIcon>

          <AiOutlineClose
            className={
              'absolute right-5 h-6 w-6 transition-all duration-300 ' +
              (!isMenuOpen ? 'opacity-0' : 'opacity-100')
            }
          />
        </MobileButton>
      </div>
    </HeaderContainer>
  )
}

export default Header

const NavButton = ({
  children,
  items,
  href,
  to,
}: {
  children: React.ReactNode
  items?: PopoverItem[]
  href?: string
  to?: string
}) => {
  return (
    <>
      {items && (
        <PopoverContainer items={items}>
          <Popover.Button className="rounded-2xl px-4 py-2 text-lg font-semibold text-gray-800 outline-none transition-all duration-150 hover:bg-backgroundGray">
            {children}
          </Popover.Button>
        </PopoverContainer>
      )}

      {href && (
        <a
          href={href || '/'}
          target="_blank"
          className="rounded-2xl px-4 py-2 text-lg font-semibold text-gray-800 outline-none transition-all duration-150 hover:bg-backgroundGray"
        >
          {children}
        </a>
      )}

      {to && (
        <Link
          to={to || '/'}
          className="rounded-2xl px-4 py-2 text-lg font-semibold text-gray-800 outline-none transition-all duration-150 hover:bg-backgroundGray"
        >
          {children}
        </Link>
      )}
    </>
  )
}
