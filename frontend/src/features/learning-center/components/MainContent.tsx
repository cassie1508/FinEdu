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
      className="flex flex-col gap-4 p-6 flex-1"
      style={{
        background: 'linear-gradient(135deg, #9EC0FF 0%, #F8FAFF 45%, #FFF6E2 72%, #FFDF94 100%)',
      }}
    >
      {/* Flashcard Section */}
      <FlashcardSection
        flashcards={flashcards}
        onCreateFlashcard={onCreateFlashcard}
        onUpdateFlashcard={onUpdateFlashcard}
        onDeleteFlashcard={onDeleteFlashcard}
      />

      {/* Resource Section */}
      <ResourceSection />
    </div>
  );
}
