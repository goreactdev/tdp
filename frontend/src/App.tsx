import './App.css'
import type { PartialColorsSet, THEME } from '@tonconnect/ui-react'
import { TonConnectUIProvider } from '@tonconnect/ui-react'
import React, { Suspense } from 'react'
import { Provider } from 'react-redux'
import {
  createBrowserRouter,
  Navigate,
  RouterProvider,
  Outlet,
} from 'react-router-dom'
import { PersistGate } from 'redux-persist/integration/react'
import { DefaultContainer } from './components/DefaultContainer'
import Footer from './components/Footer/Footer'
import Header from './components/Header/Header'
import { Loader } from './components/Loader'
import ScrollToTop from './components/ScrollToTop'
import { AppWrapper } from './components/Wrapper'
import { useMemoizedUser } from './hooks/useMemoizedUser'
import NotFoundPage from './pages/NotFound'
import { persistor, store } from './store'

const lazyReactNaiveRetry: typeof React.lazy = (importer) => {
  const retryImport = async () => {
    try {
      return await importer();
    } catch (error) {
      // retry 5 times with 1 second delay
      for (let i = 0; i < 5; i++) {
        await new Promise((resolve) => setTimeout(resolve, 1000));
        try {
          return await importer();
        } catch (e) {
          console.log("retrying import");
        }
      }
      throw error;
    }
  };
  return React.lazy(retryImport);
};

const RankingPage = lazyReactNaiveRetry(() => import('./pages/Ranking'))
const LandingPage = lazyReactNaiveRetry(() => import('./pages/Landing'))
const ProfilePage = lazyReactNaiveRetry(() => import('./pages/Profile'))
const SettingsPage = lazyReactNaiveRetry(() => import('./pages/Settings'))
const AchievementsPage = lazyReactNaiveRetry(() => import('./pages/Achievements'))

const Root = () => {
  return (

      <div className="fixed z-40 h-full w-full  overflow-hidden">
        <div id="container" className="h-full w-full  overflow-auto">
          <Header />
          <AppWrapper>
            <DefaultContainer>
              <ScrollToTop />
              <Outlet />
            </DefaultContainer>

            <Footer />
          </AppWrapper>
        </div>
      </div>
  )
}

const ProtectedRoute = ({ children }: { children: JSX.Element }) => {
  const { user } = useMemoizedUser()
  if (!user) {
    // user is not authenticated
    return <Navigate to="/" />
  }
  return children
}


const router = createBrowserRouter([
  {
    children: [
      {
        element: (
          <Suspense fallback={<Loader />}>
            <LandingPage />
          </Suspense>
        ),
        errorElement: <NotFoundPage />,
        path: '/',
      },

      {
        element: (
          <Suspense fallback={<Loader />}>
            <ProfilePage />
          </Suspense>
        ),
        errorElement: <NotFoundPage />,
        path: '/user/:username',
      },
      {
        element: (
          <Suspense fallback={<Loader />}>
            <RankingPage />
          </Suspense>
        ),
        errorElement: <NotFoundPage />,
        path: '/rewards',
      },
      {
        element: (
          <ProtectedRoute>
            <Suspense fallback={<Loader />}>
              <AchievementsPage />
            </Suspense>
          </ProtectedRoute>
        ),
        errorElement: <NotFoundPage />,
        path: '/achievements',
      },

      {
        element: (
          <ProtectedRoute>
            <Suspense fallback={<Loader />}>
              <SettingsPage />
            </Suspense>
          </ProtectedRoute>
        ),
        errorElement: <NotFoundPage />,
        path: '/settings',
      },
    ],
    element: <Root />,
    errorElement: <NotFoundPage />,
    path: '/',
  },
])


const myColorsSet: Partial<Record<THEME, PartialColorsSet>> = {
  LIGHT: {
    connectButton: {
      background: '#0088CC',
      foreground: 'white',
    },
  },
}


export const App: React.FC = () => {
  return (
    <TonConnectUIProvider
      uiPreferences={{
        colorsSet: myColorsSet,
      }}
      manifestUrl="https://tdp.tonbuilders.com/v1/manifest-ton-connect"
    >
      <Provider store={store}>
        <PersistGate loading={null} persistor={persistor}>
            <RouterProvider router={router} />
        </PersistGate>
      </Provider>
    </TonConnectUIProvider>
  )
}

export default App
