import { Field, ErrorMessage, Form } from 'formik'
import styled from 'styled-components'
import tw from 'twin.macro'

export const FormField = styled(Field)`
  ${tw` resize-none border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500
  transition-all duration-300 ease-in-out  
  
  `}
`

export const FormLabel = styled.label`
  ${tw`font-medium text-sm`}
`

export const FormError = styled(ErrorMessage)`
  ${tw`text-red-400 text-sm`}
`

export const BasicForm = styled(Form)`
  ${tw`mt-4`}
`
