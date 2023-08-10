type TableProps<T extends Record<string, string | number | React.ReactNode>> = {
  isUserRating?: boolean
  data: T[]
  config: {
    key: keyof T
    label: string
  }[]
}
const getBackgroundColor = (i: number, isUserRating: boolean): string => {
  if (isUserRating) {
    if (i % 2 === 0 && i < 10) return `bg-[#f1f7fb] bg-opacity-60`
    if (i % 1 === 0 && i < 10) return `bg-[#f1f7fb] bg-opacity-20`
    if (i % 2 === 0 && i < 20) return `bg-amber-50`
    if (i % 1 === 0 && i < 20) return `bg-amber-50 bg-opacity-60`
    if (i % 2 === 0 && i < 30) return `bg-slate-50`
    if (i % 1 === 0 && i < 30) return `bg-slate-50 bg-opacity-60`
    return `bg-backgroundGray`
  }
  return i % 2 === 0 ? `bg-backgroundGray` : `bg-white`
}

export const Table = <
  T extends Record<string, string | number | React.ReactNode>
>({
  data,
  config,
  isUserRating = false,
}: TableProps<T>) => {
  return (
    <div className="flex flex-col">
      <div className="max-w-full overflow-x-scroll rounded-2xl bg-white py-6 sm:overflow-auto">
        <table className="min-w-full divide-y divide-gray-200">
          <thead>
            <tr>
              {config.map(({ label }) => (
                <th
                  scope="col"
                  className="px-6 py-3 text-left text-xs font-semibold uppercase text-gray-500"
                >
                  {label}
                </th>
              ))}
            </tr>
          </thead>
          <tbody className={`divide-y divide-gray-200  bg-white`}>
            {data.map((row, i) => (
              <tr className={getBackgroundColor(i, isUserRating)}>
                {config.map(({ key }) => {
                  return (
                    <td className="whitespace-nowrap  px-6 py-4 text-sm text-gray-900">
                      {typeof row[key] === 'number' ? row[key] : row[key]}
                    </td>
                  )
                })}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
