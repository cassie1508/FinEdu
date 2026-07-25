import { colors } from '../lib/colors';
import { FlashcardSection } from './FlashcardSection';
import { ResourceSection } from './ResourceSection';
import { Flashcard } from '../lib/types';
import { FlashcardInput } from '../lib/flashcardsApi';

interface MainContentProps {
  flashcards: Flashcard[];
  onCreateFlashcard: (input: FlashcardInput) => Promise<void>;
  onUpdateFlashcard: (id: string, input: FlashcardInput) => Promise<void>;
  onDeleteFlashcard: (id: string) => Promise<void>;
}

export function MainContent({ flashcards, onCreateFlashcard, onUpdateFlashcard, onDeleteFlashcard }: MainContentProps) {
  return (
    <div
      className="flex overflow-hidden flex flex-col gap-2 p-6 flex-1"
      style={{ 
        background: 'linear-gradient(135deg, #9EC0FF 0%, #F8FAFF 45%, #FFF6E2 72%, #FFDF94 100%)',
        scrollbarGutter: 'stable',
      }}
    >
      {/* Flashcard Section - Left/Top with scroll */}
      <div className="flex-1 min-h-0">
        <FlashcardSection
          flashcards={flashcards}
          onCreateFlashcard={onCreateFlashcard}
          onUpdateFlashcard={onUpdateFlashcard}
          onDeleteFlashcard={onDeleteFlashcard}
        />
      </div>

      {/* Resource Section - Right/Bottom with scroll */}
      <div className="flex-1 min-h-0">
        <ResourceSection />
      </div>
    </div>
  );
}
