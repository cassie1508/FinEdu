import { BrowserRouter, Route, Routes } from 'react-router-dom'
import { Layout } from './components/layout/Layout'
import { HomePage } from './pages/HomePage'
import { CompanyDashboardPage } from './features/company-dashboard/pages/CompanyDashboardPage'
import { NewsChartsPage } from './features/news-charts/pages/NewsChartsPage'
import { LearningCenterPage } from './features/learning-center/pages/LearningCenterPage'
import { PortfolioPage } from './features/portfolio/pages/PortfolioPage'

function App() {
  return (
    <BrowserRouter>
      <div className="h-screen w-screen overflow-hidden bg-brand-bg">
        <Layout>
          <Routes>
            <Route path="/" element={<HomePage />} />
            <Route path="/dashboard" element={<CompanyDashboardPage />} />
            <Route path="/news" element={<NewsChartsPage />} />
            <Route path="/learn" element={<LearningCenterPage />} />
            <Route path="/portfolio" element={<PortfolioPage />} />
          </Routes>
        </Layout>
      </div>
    </BrowserRouter>
  )
}

export default App
