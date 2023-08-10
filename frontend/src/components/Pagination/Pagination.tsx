import React from 'react'

type PaginationProps = {
  currentPage: number
  setCurrentPage: React.Dispatch<React.SetStateAction<number>>
  pageSize: string
  count: number
}

const Pagination = ({
  currentPage,
  setCurrentPage,
  pageSize,
  count,
}: PaginationProps) => {
  const maxPages = Math.ceil(count / parseInt(pageSize)) || 1
  const [pageNumbers, setPageNumbers] = React.useState([1, 2, 3, 4, 5])

  const nextPage = () => {
    if (currentPage < maxPages) {
      setCurrentPage((prev) => prev + 1)
    }
  }

  const prevPage = () => {
    if (currentPage > 1) {
      setCurrentPage((prev) => prev - 1)
    }
  }

  React.useEffect(() => {
    let startPage = currentPage - 2
    let endPage = currentPage + 2

    // If currentPage is close to the start, show first pages
    if (currentPage < 4) {
      startPage = 1
      endPage = 5
    }

    // If currentPage is close to the end, show last pages
    if (currentPage > maxPages - 3) {
      startPage = maxPages - 4
      endPage = maxPages
    }

    // Always include 1 in page numbers
    startPage = Math.max(2, startPage)

    // Always include maxPages in page numbers
    endPage = Math.min(maxPages - 1, endPage)

    const newPageNumbers = [
      1,
      ...Array(endPage - startPage + 1)
        .fill(0)
        .map((_, idx) => startPage + idx),
      maxPages,
    ].filter((num, idx, arr) => arr.indexOf(num) === idx)

    setPageNumbers(newPageNumbers)
  }, [currentPage, maxPages])

  React.useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'ArrowRight') {
        nextPage()
      }
      if (e.key === 'ArrowLeft') {
        prevPage()
      }
    }
    window.addEventListener('keydown', handleKeyDown)
    return () => {
      window.removeEventListener('keydown', handleKeyDown)
    }
  }, [currentPage]) // Changed dependency to currentPage

  return (
    <div className="flex flex-col items-end justify-between bg-white py-3 sm:px-6">
      <div className="flex w-full flex-1 items-center justify-between sm:w-auto">
        <button
          onClick={prevPage}
          disabled={currentPage === 1}
          className={`relative inline-flex items-center rounded-2xl border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 ${
            currentPage === 1 ? 'cursor-not-allowed opacity-50' : ''
          }`}
        >
          Previous
        </button>
        <div className="hidden sm:block">
          <nav
            className="isolate inline-flex rounded-2xl px-4 "
            aria-label="Pagination"
          >
            {pageNumbers.map((number) => {
              if (number === 1)
                return (
                  <>
                    <button
                      key={number}
                      onClick={() => {
                        setCurrentPage(number)
                        setPageNumbers([1, 2, 3, 4, 5])
                      }}
                      className={`relative mx-1 inline-flex items-center rounded-xl px-4 py-2 text-sm font-medium ${
                        currentPage === number
                          ? 'z-10 bg-mainColor text-white focus:z-20 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600'
                          : 'text-gray-900 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0'
                      }`}
                    >
                      {number}
                    </button>
                    {pageNumbers[1] !== 2 && (
                      <div className="flex w-[6rem] cursor-not-allowed items-center justify-center rounded-2xl bg-gray-100 text-gray-900 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0">
                        ---
                      </div>
                    )}
                  </>
                )

              return (
                <button
                  key={number}
                  onClick={() => setCurrentPage(number)}
                  className={`relative mx-1 inline-flex items-center rounded-xl px-4 py-2 text-sm font-medium transition-all duration-500 ease-in-out ${
                    currentPage === number
                      ? 'z-10 bg-mainColor text-white focus:z-20 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600'
                      : 'text-gray-900 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0'
                  }`}
                >
                  {number}
                </button>
              )
            })}
          </nav>
        </div>
        <button
          onClick={nextPage}
          disabled={currentPage === maxPages}
          className={`relative inline-flex items-center rounded-2xl border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 ${
            currentPage === Math.ceil(maxPages)
              ? 'cursor-not-allowed opacity-50'
              : ''
          }`}
        >
          Next
        </button>

        {currentPage !== Math.ceil(maxPages) && (
          <button
            onClick={() => {
              setCurrentPage(Math.ceil(maxPages))
              setPageNumbers([
                Math.ceil(maxPages) - 4,
                Math.ceil(maxPages) - 3,
                Math.ceil(maxPages) - 2,
                Math.ceil(maxPages) - 1,
                Math.ceil(maxPages),
              ])
            }}
            disabled={currentPage === maxPages}
            className={`relative ml-2 inline-flex items-center rounded-2xl border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 ${
              currentPage === Math.ceil(maxPages)
                ? 'cursor-not-allowed opacity-50'
                : ''
            }`}
          >
            Last
          </button>
        )}
      </div>
      <div className="mt-2  ">
        <p className="text-sm text-gray-700">
          Showing
          <span className="mx-1 font-medium">
            {(currentPage - 1) * parseInt(pageSize) + 1}
          </span>
          to
          <span className="mx-1 font-medium">
            {currentPage * parseInt(pageSize)}
          </span>
          of
          <span className="mx-1 font-medium">
            {maxPages * parseInt(pageSize)}
          </span>
          results
        </p>
      </div>
    </div>
  )
}

export default Pagination
