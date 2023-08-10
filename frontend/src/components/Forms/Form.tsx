import type { FormikHelpers, FormikValues } from 'formik'
import { Formik, FieldArray } from 'formik'
import React from 'react'
import { AiFillCloseCircle } from 'react-icons/ai'
import type * as Yup from 'yup'

import Button from '../Button/Button'

import { FormError, BasicForm, FormField, FormLabel } from './Form.styles'

// generic function to render form fields

const fields = (initialValues: FormikValues) => {
  const fieldsWithoutArray = Object.keys(initialValues).filter(
    (field) => !Array.isArray(initialValues[field])
  )

  return fieldsWithoutArray.map((field, index) => (
    <div key={index} className="w-full">
      <FormLabel className="text-sm font-medium" htmlFor={field}>
        {field.replace('_', ' ')[0].toUpperCase() +
          field.replace('_', ' ').slice(1)}
      </FormLabel>
      <FormField
        placeholder="write here something..."
        type="text"
        name={field}
      />
      <FormError name={field} component="div" />
    </div>
  ))
}

const fieldArray = (
  fieldName: string,
  labelText: string,
  maxItems: number,
  values: FormikValues
) => {
  const allValues = values[fieldName]

  return (
    <div>
      <FormLabel>{labelText}</FormLabel>
      <FieldArray name={fieldName}>
        {({ push, remove }) => (
          <div className="mb-4 space-y-4">
            {allValues &&
              allValues.map((_: string, index: number) => (
                <div key={index}>
                  <div className="flex items-center">
                    <FormField name={`${fieldName}.${index}`} />
                    <AiFillCloseCircle
                      className="ml-2 cursor-pointer fill-red-400"
                      size={25}
                      type="button"
                      onClick={() => remove(index)}
                    />
                  </div>
                  <FormError name={`${fieldName}.${index}`} component="div" />
                </div>
              ))}

            {allValues && allValues.length < maxItems && (
              <Button
                type="button"
                onClick={() => push('')}
                className="mt-2"
                color="white"
              >
                Add {labelText.slice(0, -1)}
              </Button>
            )}
          </div>
        )}
      </FieldArray>
    </div>
  )
}

type ArrayField<T> = {
  fieldName: Extract<keyof T, string>
  labelText: string
  maxItems: number
  values?: T
}

type FormProps<T> =
  | {
      initialValues: T
      validationSchema: Yup.Schema<T>
      onSubmit: (
        values: T,
        formikHelpers: FormikHelpers<T>
      ) => void | Promise<void>
      arrayFields?: ArrayField<T>[]
    }
  | {
      initialValues?: undefined
      validationSchema: Yup.Schema<T>
      onSubmit: (
        values: T,
        formikHelpers: FormikHelpers<T>
      ) => void | Promise<void>
      arrayFields: ArrayField<T>[]
    }

const Form: React.FC<FormProps<FormikValues>> = ({
  initialValues,
  validationSchema,
  onSubmit,
  arrayFields,
}) => {
  // add to arrayFieldsInitialValues field.values

  const arrayFieldsValues = arrayFields?.reduce(
    (acc, field) => ({ ...acc, [field.fieldName]: field.values }),
    []
  )

  return (
    <Formik
      initialValues={{ ...initialValues, ...arrayFieldsValues }}
      validationSchema={validationSchema}
      onSubmit={onSubmit}
    >
      {({ isSubmitting, values }) => (
        <>
          <BasicForm>
            <div className="grid-cols-2 gap-4  space-y-4 sm:grid sm:space-y-0">
              {initialValues && fields(initialValues)}
              {arrayFields &&
                arrayFields.map((field) =>
                  fieldArray(
                    field.fieldName,
                    field.labelText,
                    field.maxItems,
                    values
                  )
                )}
            </div>

            <div className="mt-4 flex items-end justify-end sm:mt-0">
              <Button type="submit" disabled={isSubmitting} color="blue">
                {isSubmitting ? 'Submitting...' : 'Save Changes'}
              </Button>
            </div>
          </BasicForm>
        </>
      )}
    </Formik>
  )
}

export default Form
