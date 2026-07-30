import { useState } from 'react'
import { StockChart } from '../components/StockChart'
import { TickerNewsPanel } from '../components/TickerNewsPanel'

export function MarketDashboardPage() {
  const [symbol, setSymbol] = useState('AAPL')

  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
        Market Dashboard
      </h1>
      <p className="mt-2 text-gray-600 dark:text-gray-400">
        Explore interactive stock charts to track price movements and analyze market trends, with
        news and AI-generated insights for the company you're viewing.
      </p>

      <div className="mt-6 grid grid-cols-1 gap-6 lg:grid-cols-[minmax(0,1fr)_380px] lg:items-start">
        <StockChart symbol={symbol} onSymbolChange={setSymbol} />
        <TickerNewsPanel symbol={symbol} />
      </div>
    </div>
  )
}
