import { useEffect, useState } from 'react';
import { MessageCircle } from 'lucide-react';
import { AIChatSidebar } from '../components/AIChatSidebar';
import { MainContent } from '../components/MainContent';
import { RightSidebar } from '../components/RightSidebar';
import { colors } from '../lib/colors';
import { Navbar } from '../../../components/layout/Navbar';
import { Flashcard } from '../lib/types';
import { createFlashcard, deleteFlashcard, fetchFlashcards, FlashcardInput, updateFlashcard } from '../lib/flashcardsApi';

export function LearningCenterPage() {
  const [isChatOpen, setIsChatOpen] = useState(false);
  const [flashcards, setFlashcards] = useState<Flashcard[]>([]);
  const [flashcardsError, setFlashcardsError] = useState<string | null>(null);

  useEffect(() => {
    fetchFlashcards()
      .then(setFlashcards)
      .catch(() => setFlashcardsError('Failed to load flashcards from the server.'));
  }, []);

  const handleCreateFlashcard = async (input: FlashcardInput) => {
    const created = await createFlashcard(input);
    setFlashcards(prev => [...prev, created]);
  };

  const handleUpdateFlashcard = async (id: string, input: FlashcardInput) => {
    const updated = await updateFlashcard(id, input);
    setFlashcards(prev => prev.map(card => (card.id === id ? updated : card)));
  };

  const handleDeleteFlashcard = async (id: string) => {
    await deleteFlashcard(id);
    setFlashcards(prev => prev.filter(card => card.id !== id));
  };

  return (
    <div className="relative min-h-screen">
      <Navbar />
      {flashcardsError && (
        <p className="px-6 pt-2 text-sm" style={{ color: colors.accent }}>
          {flashcardsError}
        </p>
      )}
      <div className="flex">
        {/* Main Content Area */}
        <MainContent
          flashcards={flashcards}
          onCreateFlashcard={handleCreateFlashcard}
          onUpdateFlashcard={handleUpdateFlashcard}
          onDeleteFlashcard={handleDeleteFlashcard}
        />

        {/* Right Sidebar - Analytics */}
        <aside className="w-80 flex-shrink-0">
          <RightSidebar />
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
