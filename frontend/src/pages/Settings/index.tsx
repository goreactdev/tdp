import type { ChangeEvent } from 'react'
import { useEffect, useRef } from 'react'
import { AiFillGithub } from 'react-icons/ai'
import { BsTelegram } from 'react-icons/bs'
import { FaUserAlt } from 'react-icons/fa'
import * as Yup from 'yup'

import Button from '../../components/Button'
import Form from '../../components/Forms'
import { jobs } from '../../components/Landing/MeetWinners/MeetWinners'
import { PageHeaderText } from '../../components/Text'
import type { LinkedAccount } from '../../services/types'
import {
  useCheckAuthTelegramMutation,
  useGetMyAccountQuery,
  useUnlinkAccountMutation,
  useUpdateUserMutation,
  useUploadImageMutation,
} from '../../services/users'
import { BASE_URL, BOT_ID } from '../../utils/config'

const validationSchema = Yup.object().shape({
  bio: Yup.string()
    .min(3, 'Bio must be at least 3 characters')
    .max(500, 'Bio must be less than 500 characters'),
  certifications: Yup.array().of(
    Yup.string().max(50, 'Certification must be less than 50 characters')
  ),
  first_name: Yup.string()
    .min(3, 'First name must be at least 3 characters')
    .max(50, 'First name must be less than 50 characters'),
  job: Yup.string()
    .min(3, 'Job must be at least 3 characters')
    .max(100, 'Job must be less than 50 characters'),
  languages: Yup.array().of(
    Yup.string().max(50, 'Language must be less than 50 characters')
  ),
  last_name: Yup.string()
    .min(3, 'Last name must be at least 3 characters')
    .max(50, 'Last name must be less than 50 characters'),
  username: Yup.string()
    .min(3, 'Username must be at least 3 characters')
    .max(50, 'Username must be less than 20 characters')
    .matches(
      // could be uppercase letters
      /^[a-zA-Z0-9-_]+$/,
      'Username can only contain lowercase letters, numbers, hyphens, and underscores'
    ),
})

type AccountSectionProps = {
  accountType: 'github' | 'telegram'
  icon: React.ReactNode
  linked_accounts: LinkedAccount[] | undefined
  handleConnect: () => void
  unlinkAccount: (provider: { provider: string }) => void
}

const AccountSection = ({
  accountType,
  icon,
  linked_accounts,
  handleConnect,
  unlinkAccount,
}: AccountSectionProps) => {
  const account = linked_accounts?.find((acc) => acc.provider === accountType)

  return (
    <div className="flex w-full items-center justify-between py-4">
      <div className="flex items-center space-x-3">
        {icon}

        <div className="leading-6">
          <h3 className="text-xs font-semibold sm:text-base">
            {accountType.charAt(0).toUpperCase() + accountType.slice(1)} Account
          </h3>

          {account ? (
            <a
              href={`https://${
                accountType === 'telegram' ? 't.me' : 'github.com'
              }/${account.login}`}
              target="_blank"
              className="hidden text-sm text-blue-600 sm:block"
            >
              {`https://${accountType === 'telegram' ? 't.me' : 'github.com'}/${
                account.login
              }`}
            </a>
          ) : (
            <p className="hidden text-sm text-gray-600 sm:block">
              Not Connected
            </p>
          )}
        </div>
      </div>

      {account ? (
        <Button
          color="white"
          onClick={() => {
            unlinkAccount({ provider: accountType })
            // if it's a telegram account, disconnect it
            if (accountType === 'telegram') {
              // delete the hash from the url
              window.location.hash = ''
            }
          }}
        >
          Disconnect
        </Button>
      ) : (
        <Button color="blue" onClick={handleConnect}>
          Connect
        </Button>
      )}
    </div>
  )
}

const ImageUploader = ({ children }: { children: React.ReactNode }) => {
  // const [updateUser] = useUpdateUserMutation();

  const inputRef = useRef<HTMLInputElement>(null)

  const handleClick = () => {
    if (inputRef.current) inputRef.current.click()
  }

  const [uploadImage] = useUploadImageMutation()

  const [updateUser] = useUpdateUserMutation()

  const onFileChange = async (e: ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      const selectedFile = e.target.files[0]
      await onFileUpload(selectedFile)
    }
  }

  const onFileUpload = async (file: File) => {
    const formData = new FormData()
    formData.append('file', file)

    const response = await uploadImage(formData).unwrap()

    // The URL of the image is assumed to be in the 'url' field of the response
    updateUser({ avatar_url: response.url })
  }

  return (
    <label className="mt-4" onClick={handleClick}>
      <input
        ref={inputRef}
        className="hidden"
        type="file"
        hidden
        onChange={onFileChange}
      />
      {children}
    </label>
  )
}

const SettingsPage = () => {
  const [trigger] = useUpdateUserMutation()

  const [unlinkAccount] = useUnlinkAccountMutation()

  const [checkAuthTelegram] = useCheckAuthTelegramMutation()

  const handleConnectGithub = () => {
    // This should be the URL of your Go server's OAuth handler
    const oauthURL = BASE_URL + '/v1/github/login?username=' + data?.username
    window.location.href = oauthURL
  }

  // get #tgAuthResult from url
  const hash = window.location.hash.substring(1)
  const tgAuthResult = hash.split('=')[1]

  const handleConnectTelegram = () => {
    const authTelegram = `https://oauth.telegram.org/auth?bot_id=${BOT_ID}&scope=bot&origin=${encodeURIComponent(BASE_URL + "/settings/")}&request_access=write`

    window.location.href = authTelegram
  }

  const { data } = useGetMyAccountQuery(undefined)

  useEffect(() => {
    if (tgAuthResult) {
      checkAuthTelegram({ auth_obj: tgAuthResult })
    }
  }, [tgAuthResult])

  return (
    <>
      <PageHeaderText>User settings</PageHeaderText>

      {data && (
        <div className="mt-2 grid-cols-8 gap-4 space-y-4 lg:mt-6 lg:grid lg:space-y-0 ">
          <div className="col-span-4 h-max">
            <div className="flex flex-col items-center rounded-2xl bg-gray-50 p-4 sm:flex-row sm:items-start sm:p-6">
              {data?.avatar_url ? (
                <img
                  src={data?.avatar_url ?? '/images/default-avatar.png'}
                  alt=""
                  className="h-32 w-32 rounded-full object-cover object-center"
                />
              ) : (
                <div className="flex h-32 w-32 items-center justify-center rounded-full bg-mainColor transition-transform duration-300 hover:scale-105">
                  <FaUserAlt size={90} className="p-2" fill="white" />
                </div>
              )}

              <div className="mt-2 flex flex-col justify-between text-center sm:ml-4 sm:text-start">
                <div className="flex flex-col space-y-1">
                  <span className="text-xl font-bold text-mainColor">
                    {data?.first_name + ' ' + data?.last_name}
                  </span>
                  <span className="text-base font-medium text-gray-600">
                    {data?.job ||
                      jobs[data.id > 9 ? Number(String(data.id)[0]) : data.id]}
                  </span>
                </div>
                <ImageUploader>
                  {/* <div className="mt-4 sm:mt-0"> */}
                  <Button color="blue">Change Picture</Button>
                  {/* </div> */}
                </ImageUploader>
              </div>
            </div>
            <div className="col-span-4 mt-4 h-max rounded-2xl bg-gray-50 p-6">
              <h2 className="flex items-center text-xl font-bold">
                Linked Accounts{' '}
                <span className="ml-2 text-sm text-mainColor">+200 points</span>
              </h2>

              <div className="divide-y">
                {data && (
                  <>
                    <AccountSection
                      accountType="github"
                      icon={<AiFillGithub className="h-8 w-8" />}
                      linked_accounts={data.linked_accounts}
                      handleConnect={handleConnectGithub}
                      unlinkAccount={unlinkAccount}
                    />

                    <AccountSection
                      accountType="telegram"
                      icon={<BsTelegram className="h-8 w-8" />}
                      linked_accounts={data.linked_accounts}
                      handleConnect={handleConnectTelegram}
                      unlinkAccount={unlinkAccount}
                    />
                  </>
                )}
              </div>
            </div>

            <div className="mt-4 rounded-2xl bg-gray-50 p-4 sm:p-6">
              <h2 className="text-xl font-bold">Additional Information</h2>
              <Form
                validationSchema={validationSchema}
                onSubmit={(values, { setSubmitting }) => {
                  trigger({
                    certifications: values.certifications,
                    languages: values.languages,
                  })

                  setTimeout(() => {
                    setSubmitting(false)
                  }, 400)
                }}
                arrayFields={[
                  {
                    fieldName: 'languages',
                    labelText: 'Languages',
                    maxItems: 5,
                    values: data?.languages || [],
                  },
                  {
                    fieldName: 'certifications',
                    labelText: 'Certifications',
                    maxItems: 5,
                    values: data?.certifications || [],
                  },
                ]}
              />
            </div>
          </div>

          <div className="col-span-4 h-max rounded-2xl bg-gray-50 p-6">
            <h2 className="text-xl font-bold">General Information</h2>
            <Form
              validationSchema={validationSchema}
              onSubmit={(values, { setSubmitting }) => {
                trigger({
                  bio: values.bio,
                  certifications: values.certifications,
                  first_name: values.first_name,
                  job: values.job,
                  languages: values.languages,
                  last_name: values.last_name,
                  username: values.username,
                })

                setTimeout(() => {
                  setSubmitting(false)
                }, 400)
              }}
              initialValues={{
                bio: data?.bio || '',
                first_name: data?.first_name,
                job: data?.job || '',
                last_name: data?.last_name,
                username: data?.username,
              }}
            />
          </div>
        </div>
      )}
    </>
  )
}

export default SettingsPage
