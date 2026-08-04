import { lazy, Suspense } from 'react'
import { BrowserRouter, Route, Routes } from 'react-router-dom'
import { Layout } from './components/layout/Layout'
import { HomePage } from './pages/HomePage'
import { MarketDashboardPage } from './features/market-dashboard/pages/MarketDashboardPage'
import { NewsPage } from './features/news/pages/NewsPage'
import { LearningCenterPage } from './features/learning-center/pages/LearningCenterPage'
import { LoginPage } from './features/auth/pages/LoginPage'
import { SignUpPage } from './features/auth/pages/SignUpPage'
import { RequireGuest } from './features/auth/components/RequireGuest'
import { RequireAuth } from './features/auth/components/RequireAuth'

const PortfolioPage = lazy(() =>
  import('./features/portfolio/pages/PortfolioPage').then((module) => ({
    default: module.PortfolioPage,
  })),
)

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route
          path="/login"
          element={
            <RequireGuest>
              <LoginPage />
            </RequireGuest>
          }
        />
        <Route
          path="/register"
          element={
            <RequireGuest>
              <SignUpPage />
            </RequireGuest>
          }
        />
        <Route
          path="/*"
          element={
            <Layout>
              <Routes>
                <Route path="/" element={<HomePage />} />
                <Route path="/dashboard" element={<MarketDashboardPage />} />
                <Route path="/news" element={<NewsPage />} />
                <Route
                  path="/learn"
                  element={
                    <RequireAuth>
                      <LearningCenterPage />
                    </RequireAuth>
                  }
                />
                <Route
                  path="/portfolio"
                  element={
                    <RequireAuth>
                      <Suspense fallback={<div className="h-80 animate-pulse rounded-2xl bg-[#E3DEDE]" />}>
                        <PortfolioPage />
                      </Suspense>
                    </RequireAuth>
                  }
                />
              </Routes>
            </Layout>
          }
        />
      </Routes>
    </BrowserRouter>
  )
}

export default App
