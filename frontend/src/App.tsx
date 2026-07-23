import { lazy, Suspense } from 'react'
import { BrowserRouter, Route, Routes } from 'react-router-dom'
import { Layout } from './components/layout/Layout'
import { HomePage } from './pages/HomePage'
import { CompanyDashboardPage } from './features/company-dashboard/pages/CompanyDashboardPage'
import { NewsChartsPage } from './features/news-charts/pages/NewsChartsPage'
import { LearningCenterPage } from './features/learning-center/pages/LearningCenterPage'

const PortfolioPage = lazy(() =>
  import('./features/portfolio/pages/PortfolioPage').then((module) => ({
    default: module.PortfolioPage,
  })),
)

function App() {
  return (
    <BrowserRouter>
      <Layout>
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/dashboard" element={<CompanyDashboardPage />} />
          <Route path="/news" element={<NewsChartsPage />} />
          <Route path="/learn" element={<LearningCenterPage />} />
          <Route
            path="/portfolio"
            element={(
              <Suspense fallback={<div className="h-80 animate-pulse rounded-2xl bg-[#E3DEDE]" />}>
                <PortfolioPage />
              </Suspense>
            )}
          />
        </Routes>
      </Layout>
    </BrowserRouter>
  )
}

export default App
