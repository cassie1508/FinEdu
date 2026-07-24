import { useState } from 'react';
import { MessageCircle } from 'lucide-react';
import { AIChatSidebar } from '../components/AIChatSidebar';
import { MainContent } from '../components/MainContent';
import { RightSidebar } from '../components/RightSidebar';
import { mockFlashcards, mockProgress } from '../lib/mockData';
import { colors } from '../lib/colors';
import { Navbar } from '../../../components/layout/Navbar';

export function LearningCenterPage() {
  const [isChatOpen, setIsChatOpen] = useState(false);

  return (
    <div className="relative h-screen overflow-hidden" style={{ backgroundColor: colors.background }}>
      <Navbar />
      <div className="flex h-full overflow-hidden">
        {/* Main Content Area */}
        <MainContent flashcards={mockFlashcards} />

        {/* Right Sidebar - Analytics */}
        <aside className="w-80 flex-shrink-0 overflow-y-scroll scrollbar-custom">
        <RightSidebar progress={mockProgress} />
      </aside>
      </div>

      {/* Floating chat trigger */}
      <button
        onClick={() => setIsChatOpen(prev => !prev)}
        className="fixed bottom-6 right-6 z-40 w-14 h-14 rounded-full shadow-lg transition-transform hover:scale-105 flex items-center justify-center"
        style={{
          backgroundColor: colors.primary,
          color: colors.text.light,
        }}
        aria-label="Open AI chat"
      >
        <MessageCircle size={24} />
      </button>

      {/* Chat popup */}
      {isChatOpen && (
        <div className="fixed inset-0 z-50 flex items-end justify-end p-4 sm:p-6">
          <button
            className="absolute inset-0"
            style={{ backgroundColor: 'rgba(0, 0, 0, 0.25)' }}
            onClick={() => setIsChatOpen(false)}
            aria-label="Close AI chat backdrop"
          />
          <div className="relative">
            <AIChatSidebar mode="popup" onClose={() => setIsChatOpen(false)} />
          </div>
        </div>
      )}
    </div>
  );
}
