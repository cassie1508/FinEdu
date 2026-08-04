import { useState } from 'react';
import { Plus, X, Edit2, Save, Trash2, AlertTriangle } from 'lucide-react';
import { Textfit } from 'react-textfit';
import { colors } from '../lib/colors';
import { Flashcard } from '../lib/types';
import { FlashcardInput } from '../lib/flashcardsApi';
import { Pagination } from './Pagination';

const FLASHCARDS_PAGE_SIZE = 6;

interface FlashcardSectionProps {
  flashcards: Flashcard[];
  onCreateFlashcard: (input: FlashcardInput) => Promise<void>;
  onUpdateFlashcard: (id: string, input: FlashcardInput) => Promise<void>;
  onDeleteFlashcard: (id: string) => Promise<void>;
}

const categories = [
  'All',
  'Investing',
  'Budgeting',
  'Credit & Debt',
  'Markets',
  'Retirement',
  'Taxes',
];

interface FlashcardCardProps {
  card: Flashcard;
  onCardClick: (card: Flashcard) => void;
  onDelete: (card: Flashcard) => void;
}

function FlashcardCard({ card, onCardClick, onDelete }: FlashcardCardProps) {
  const handleClick = (e: React.MouseEvent) => {
    e.stopPropagation();
    e.preventDefault();
    onCardClick(card);
  };

  const handleDeleteClick = (e: React.MouseEvent) => {
    e.stopPropagation();
    e.preventDefault();
    onDelete(card);
  };

  return (
    <div
      onClick={handleClick}
      className="group relative h-48 cursor-pointer perspective transition-all duration-300 hover:shadow-lg hover:-translate-y-1"
      style={{
        backgroundColor: colors.surface,
        border: `1px solid ${colors.border}`,
        borderRadius: '18px',
        padding: '16px',
      }}
    >
      <button
        onClick={handleDeleteClick}
        className="absolute top-2 right-2 z-10 p-1.5 rounded-full opacity-0 group-hover:opacity-100 transition-opacity hover:shadow-md"
        style={{ backgroundColor: colors.secondary, color: colors.accent }}
        title="Delete flashcard"
        aria-label="Delete flashcard"
      >
        <Trash2 size={16} />
      </button>
      <div className="relative h-full flex flex-col" style={{ gap: '12px' }}>
        {/* Category Tag */}
        <div
          className="inline-flex w-fit px-3 py-1 rounded-full text-xs font-medium flex-shrink-0"
          style={{
            backgroundColor: colors.secondary,
            color: colors.text.secondary,
          }}
        >
          {card.category}
        </div>

        {/* Content - Equal spacing */}
        <div className="flex-1 flex flex-col overflow-hidden" style={{ gap: '12px' }}>
          <div className="flex-1 min-h-0">
            <Textfit mode="multi" max={20} min={14} style={{ color: colors.emphasis }} className="font-serif font-bold">
              {card.title}
            </Textfit>
          </div>
          <div className="flex-1 min-h-0">
            <Textfit mode="multi" max={12} min={9} style={{ color: colors.text.secondary }}>
              {card.definition}
            </Textfit>
          </div>
        </div>

        {/* Click to view indicator */}
        <p
          className="text-xs text-center flex-shrink-0"
          style={{ color: colors.accent, opacity: 0.6 }}
        >
          Click to view
        </p>
      </div>
    </div>
  );
}





interface FlashcardDetailModalProps {
  card: Flashcard | null;
  isOpen: boolean;
  onClose: () => void;
  onSave: (updatedCard: Flashcard) => void;
}

function FlashcardDetailModal({ card, isOpen, onClose, onSave }: FlashcardDetailModalProps) {
  const [isEditing, setIsEditing] = useState(false);
  const [editedCard, setEditedCard] = useState<Flashcard | null>(null);

  if (!isOpen || !card) return null;

  const handleEdit = () => {
    setIsEditing(true);
    setEditedCard({ ...card });
  };

  const handleSave = () => {
    if (editedCard) {
      onSave(editedCard);
      setIsEditing(false);
    }
  };

  const handleCancel = () => {
    setIsEditing(false);
    setEditedCard(null);
  };

  const handleInputChange = (field: keyof Flashcard, value: string) => {
    if (editedCard) {
      setEditedCard({ ...editedCard, [field]: value });
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center" style={{ backgroundColor: 'rgba(0, 0, 0, 0.5)' }}>
      <div
        className="bg-white rounded-lg shadow-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto"
        style={{ backgroundColor: colors.surface }}
      >
        {/* Header */}
        <div
          className="flex items-center justify-between p-6 border-b sticky top-0"
          style={{ borderColor: colors.border, backgroundColor: colors.surface }}
        >
          <h2 className="text-2xl font-serif font-bold" style={{ color: colors.emphasis }}>
            {isEditing ? 'Edit Flashcard' : 'Flashcard Details'}
          </h2>
          <div className="flex gap-2">
            {!isEditing && (
              <button
                onClick={handleEdit}
                className="p-2 rounded-lg transition-all hover:shadow-md"
                style={{ backgroundColor: colors.secondary, color: colors.primary }}
                title="Edit flashcard"
              >
                <Edit2 size={20} />
              </button>
            )}
            <button
              onClick={isEditing ? handleCancel : onClose}
              className="p-2 rounded-lg transition-all hover:shadow-md"
              style={{ backgroundColor: colors.secondary, color: colors.text.secondary }}
              title={isEditing ? 'Cancel' : 'Close'}
            >
              <X size={20} />
            </button>
          </div>
        </div>

        {/* Content */}
        <div className="p-6 space-y-6">
          {/* Category Tag */}
          <div
            className="inline-flex px-3 py-1 rounded-full text-xs font-medium"
            style={{
              backgroundColor: colors.secondary,
              color: colors.text.secondary,
            }}
          >
            {isEditing ? editedCard?.category : card.category}
          </div>

          {/* Title */}
          <div>
            <h3 className="text-xs font-semibold mb-2" style={{ color: colors.text.secondary }}>
              TITLE
            </h3>
            {isEditing ? (
              <input
                type="text"
                value={editedCard?.title || ''}
                onChange={(e) => handleInputChange('title', e.target.value)}
                className="w-full p-2 border rounded"
                style={{
                  borderColor: colors.border,
                  backgroundColor: colors.background,
                  color: colors.text.primary,
                }}
              />
            ) : (
              <p style={{ color: colors.emphasis }} className="text-lg font-bold">
                {card.title}
              </p>
            )}
          </div>

          {/* Definition */}
          <div>
            <h3 className="text-xs font-semibold mb-2" style={{ color: colors.text.secondary }}>
              DEFINITION
            </h3>
            {isEditing ? (
              <textarea
                value={editedCard?.definition || ''}
                onChange={(e) => handleInputChange('definition', e.target.value)}
                className="w-full p-2 border rounded min-h-24 resize-none"
                style={{
                  borderColor: colors.border,
                  backgroundColor: colors.background,
                  color: colors.text.primary,
                }}
              />
            ) : (
              <p style={{ color: colors.text.primary }}>{card.definition}</p>
            )}
          </div>

          {/* Why It Matters */}
          <div>
            <h3 className="text-xs font-semibold mb-2" style={{ color: colors.text.secondary }}>
              WHY IT MATTERS
            </h3>
            {isEditing ? (
              <textarea
                value={editedCard?.whyItMatters || ''}
                onChange={(e) => handleInputChange('whyItMatters', e.target.value)}
                className="w-full p-2 border rounded min-h-24 resize-none"
                style={{
                  borderColor: colors.border,
                  backgroundColor: colors.background,
                  color: colors.text.primary,
                }}
              />
            ) : (
              <p style={{ color: colors.text.primary }}>{card.whyItMatters}</p>
            )}
          </div>

          {/* Example */}
          <div>
            <h3 className="text-xs font-semibold mb-2" style={{ color: colors.text.secondary }}>
              EXAMPLE
            </h3>
            {isEditing ? (
              <textarea
                value={editedCard?.example || ''}
                onChange={(e) => handleInputChange('example', e.target.value)}
                className="w-full p-2 border rounded min-h-24 resize-none"
                style={{
                  borderColor: colors.border,
                  backgroundColor: colors.background,
                  color: colors.text.primary,
                }}
              />
            ) : (
              <p style={{ color: colors.text.primary }}>{card.example}</p>
            )}
          </div>

          {/* Common Misconceptions */}
          <div>
            <h3 className="text-xs font-semibold mb-2" style={{ color: colors.text.secondary }}>
              COMMON MISCONCEPTIONS
            </h3>
            {isEditing ? (
              <textarea
                value={editedCard?.commonMisconceptions || ''}
                onChange={(e) => handleInputChange('commonMisconceptions', e.target.value)}
                className="w-full p-2 border rounded min-h-24 resize-none"
                style={{
                  borderColor: colors.border,
                  backgroundColor: colors.background,
                  color: colors.text.primary,
                }}
              />
            ) : (
              <p style={{ color: colors.text.primary }}>{card.commonMisconceptions || 'N/A'}</p>
            )}
          </div>
        </div>

        {/* Footer */}
        {isEditing && (
          <div
            className="flex gap-3 p-6 border-t justify-end sticky bottom-0"
            style={{ borderColor: colors.border, backgroundColor: colors.surface }}
          >
            <button
              onClick={handleCancel}
              className="px-4 py-2 rounded-lg font-medium transition-all"
              style={{
                backgroundColor: colors.secondary,
                color: colors.text.secondary,
              }}
            >
              Cancel
            </button>
            <button
              onClick={handleSave}
              className="flex items-center gap-2 px-4 py-2 rounded-lg font-medium transition-all text-white"
              style={{ backgroundColor: colors.primary }}
            >
              <Save size={18} />
              Save Changes
            </button>
          </div>
        )}
      </div>
    </div>
  );
}

interface ConfirmDeleteModalProps {
  card: Flashcard | null;
  isDeleting: boolean;
  onConfirm: () => void;
  onCancel: () => void;
}

function ConfirmDeleteModal({ card, isDeleting, onConfirm, onCancel }: ConfirmDeleteModalProps) {
  if (!card) return null;

  return (
    <div
      className="fixed inset-0 z-60 flex items-center justify-center p-4"
      style={{ backgroundColor: 'rgba(0, 0, 0, 0.5)' }}
      onClick={onCancel}
    >
      <div
        role="alertdialog"
        aria-labelledby="confirm-delete-title"
        aria-describedby="confirm-delete-body"
        className="w-full max-w-sm rounded-lg shadow-xl"
        style={{ backgroundColor: colors.surface }}
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex flex-col items-center gap-3 p-6 text-center">
          <div
            className="flex h-12 w-12 items-center justify-center rounded-full"
            style={{ backgroundColor: 'rgba(168, 124, 124, 0.15)', color: colors.error }}
          >
            <AlertTriangle size={24} />
          </div>
          <h2 id="confirm-delete-title" className="text-lg font-serif font-bold" style={{ color: colors.emphasis }}>
            Delete this flashcard?
          </h2>
          <p id="confirm-delete-body" className="text-sm" style={{ color: colors.text.secondary }}>
            <span className="font-semibold" style={{ color: colors.emphasis }}>
              "{card.title}"
            </span>{' '}
            will be permanently deleted. This action cannot be undone.
          </p>
        </div>

        <div
          className="flex gap-3 p-4 border-t justify-end"
          style={{ borderColor: colors.border }}
        >
          <button
            onClick={onCancel}
            disabled={isDeleting}
            className="px-4 py-2 rounded-lg font-medium transition-all disabled:opacity-50"
            style={{
              backgroundColor: colors.secondary,
              color: colors.text.secondary,
            }}
          >
            Cancel
          </button>
          <button
            onClick={onConfirm}
            disabled={isDeleting}
            className="flex items-center gap-2 px-4 py-2 rounded-lg font-medium text-white transition-all disabled:opacity-50"
            style={{ backgroundColor: colors.error }}
          >
            <Trash2 size={16} />
            {isDeleting ? 'Deleting…' : 'Delete'}
          </button>
        </div>
      </div>
    </div>
  );
}

interface CreateFlashcardModalProps {
  isOpen: boolean;
  onClose: () => void;
  onCreate: (input: FlashcardInput) => Promise<void>;
}

function CreateFlashcardModal({ isOpen, onClose, onCreate }: CreateFlashcardModalProps) {
  const [newCard, setNewCard] = useState<Flashcard>({
    id: '',
    title: '',
    category: categories[1], // Default to first category (Investing)
    definition: '',
    example: '',
    whyItMatters: '',
    commonMisconceptions: '',
    reviewCount: 0,
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  });
  const [isSaving, setIsSaving] = useState(false);

  if (!isOpen) return null;

  const handleInputChange = (field: keyof Flashcard, value: string) => {
    setNewCard({ ...newCard, [field]: value });
  };

  const handleCreate = async () => {
    setIsSaving(true);
    try {
      await onCreate({
        title: newCard.title,
        category: newCard.category,
        whyItMatters: newCard.whyItMatters,
        definition: newCard.definition,
        example: newCard.example,
        commonMisconceptions: newCard.commonMisconceptions,
      });
      // Reset form
      setNewCard({
        id: '',
        title: '',
        category: categories[1],
        definition: '',
        example: '',
        whyItMatters: '',
        commonMisconceptions: '',
        reviewCount: 0,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      });
      onClose();
    } finally {
      setIsSaving(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center" style={{ backgroundColor: 'rgba(0, 0, 0, 0.5)' }}>
      <div
        className="bg-white rounded-lg shadow-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto"
        style={{ backgroundColor: colors.surface }}
      >
        {/* Header */}
        <div
          className="flex items-center justify-between p-6 border-b sticky top-0"
          style={{ borderColor: colors.border, backgroundColor: colors.surface }}
        >
          <h2 className="text-2xl font-serif font-bold" style={{ color: colors.emphasis }}>
            Create New Flashcard
          </h2>
          <button
            onClick={onClose}
            className="p-2 rounded-lg transition-all hover:shadow-md"
            style={{ backgroundColor: colors.secondary, color: colors.text.secondary }}
            title="Close"
          >
            <X size={20} />
          </button>
        </div>

        {/* Content */}
        <div className="p-6 space-y-6">
          {/* Category Selection */}
          <div>
            <h3 className="text-xs font-semibold mb-2" style={{ color: colors.accent }}>
              CATEGORY
            </h3>
            <select
              value={newCard.category}
              onChange={(e) => handleInputChange('category', e.target.value)}
              className="w-full p-2 border rounded"
              style={{
                borderColor: colors.border,
                backgroundColor: colors.background,
                color: colors.text.primary,
              }}
            >
              {categories.slice(1).map(cat => (
                <option key={cat} value={cat}>
                  {cat}
                </option>
              ))}
            </select>
          </div>

          {/* Title */}
          <div>
            <h3 className="text-xs font-semibold mb-2" style={{ color: colors.accent }}>
              TITLE
            </h3>
            <input
              type="text"
              value={newCard.title}
              onChange={(e) => handleInputChange('title', e.target.value)}
              placeholder="Enter flashcard title"
              className="w-full p-2 border rounded"
              style={{
                borderColor: colors.border,
                backgroundColor: colors.background,
                color: colors.text.primary,
              }}
            />
          </div>

          {/* Definition */}
          <div>
            <h3 className="text-xs font-semibold mb-2" style={{ color: colors.accent }}>
              DEFINITION
            </h3>
            <textarea
              value={newCard.definition}
              onChange={(e) => handleInputChange('definition', e.target.value)}
              placeholder="Enter the definition"
              className="w-full p-2 border rounded min-h-24 resize-none"
              style={{
                borderColor: colors.border,
                backgroundColor: colors.background,
                color: colors.text.primary,
              }}
            />
          </div>

          {/* Why It Matters */}
          <div>
            <h3 className="text-xs font-semibold mb-2" style={{ color: colors.accent }}>
              WHY IT MATTERS
            </h3>
            <textarea
              value={newCard.whyItMatters}
              onChange={(e) => handleInputChange('whyItMatters', e.target.value)}
              placeholder="Explain why this matters"
              className="w-full p-2 border rounded min-h-24 resize-none"
              style={{
                borderColor: colors.border,
                backgroundColor: colors.background,
                color: colors.text.primary,
              }}
            />
          </div>

          {/* Example */}
          <div>
            <h3 className="text-xs font-semibold mb-2" style={{ color: colors.accent }}>
              EXAMPLE
            </h3>
            <textarea
              value={newCard.example}
              onChange={(e) => handleInputChange('example', e.target.value)}
              placeholder="Provide an example"
              className="w-full p-2 border rounded min-h-24 resize-none"
              style={{
                borderColor: colors.border,
                backgroundColor: colors.background,
                color: colors.text.primary,
              }}
            />
          </div>

          {/* Common Misconceptions */}
          <div>
            <h3 className="text-xs font-semibold mb-2" style={{ color: colors.accent }}>
              COMMON MISCONCEPTIONS
            </h3>
            <textarea
              value={newCard.commonMisconceptions}
              onChange={(e) => handleInputChange('commonMisconceptions', e.target.value)}
              placeholder="What misconceptions exist about this?"
              className="w-full p-2 border rounded min-h-24 resize-none"
              style={{
                borderColor: colors.border,
                backgroundColor: colors.background,
                color: colors.text.primary,
              }}
            />
          </div>
        </div>

        {/* Footer */}
        <div
          className="flex gap-3 p-6 border-t justify-end sticky bottom-0"
          style={{ borderColor: colors.border, backgroundColor: colors.surface }}
        >
          <button
            onClick={onClose}
            className="px-4 py-2 rounded-lg font-medium transition-all"
            style={{
              backgroundColor: colors.secondary,
              color: colors.text.secondary,
            }}
          >
            Cancel
          </button>
          <button
            onClick={handleCreate}
            disabled={isSaving}
            className="flex items-center gap-2 px-4 py-2 rounded-lg font-medium transition-all text-white disabled:opacity-60"
            style={{ backgroundColor: colors.primary }}
          >
            <Plus size={18} />
            {isSaving ? 'Creating...' : 'Create Flashcard'}
          </button>
        </div>
      </div>
    </div>
  );
}

export function FlashcardSection({ flashcards, onCreateFlashcard, onUpdateFlashcard, onDeleteFlashcard }: FlashcardSectionProps) {
  const [selectedCategory, setSelectedCategory] = useState('All');
  const [selectedCard, setSelectedCard] = useState<Flashcard | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [pendingDeleteCard, setPendingDeleteCard] = useState<Flashcard | null>(null);
  const [isDeleting, setIsDeleting] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);

  const handleCategoryChange = (category: string) => {
    setSelectedCategory(category);
    setCurrentPage(1);
  };

  const handleCardClick = (card: Flashcard) => {
    setSelectedCard(card);
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    setSelectedCard(null);
  };

  const handleSaveCard = async (updatedCard: Flashcard) => {
    await onUpdateFlashcard(updatedCard.id, {
      title: updatedCard.title,
      category: updatedCard.category,
      whyItMatters: updatedCard.whyItMatters,
      definition: updatedCard.definition,
      example: updatedCard.example,
      commonMisconceptions: updatedCard.commonMisconceptions,
    });
    setIsModalOpen(false);
    setSelectedCard(null);
  };

  const handleAddButton = () => {
    setIsCreateModalOpen(true);
  };

  const handleCreateCard = async (input: FlashcardInput) => {
    await onCreateFlashcard(input);
    setIsCreateModalOpen(false);
  };

  const handleRequestDelete = (card: Flashcard) => {
    setPendingDeleteCard(card);
  };

  const handleCancelDelete = () => {
    if (isDeleting) return;
    setPendingDeleteCard(null);
  };

  const handleConfirmDelete = async () => {
    if (!pendingDeleteCard) return;
    const id = pendingDeleteCard.id;
    setIsDeleting(true);
    try {
      await onDeleteFlashcard(id);
      if (selectedCard?.id === id) {
        setIsModalOpen(false);
        setSelectedCard(null);
      }
      setPendingDeleteCard(null);
    } finally {
      setIsDeleting(false);
    }
  };

  const filteredCards =
    selectedCategory === 'All'
      ? flashcards
      : flashcards.filter(card => card.category === selectedCategory);

  const totalPages = Math.max(1, Math.ceil(filteredCards.length / FLASHCARDS_PAGE_SIZE));
  const safePage = Math.min(currentPage, totalPages);
  const pagedCards = filteredCards.slice(
    (safePage - 1) * FLASHCARDS_PAGE_SIZE,
    safePage * FLASHCARDS_PAGE_SIZE,
  );

  return (
    <div
      className="flex flex-col rounded-xl border"
      style={{
        backgroundColor: colors.surface,
        borderColor: colors.border,
      }}
    >
      {/* Header - Fixed */}
      <div className="flex-shrink-0 p-5 border-b" style={{ borderColor: colors.border }}>
        <div className="flex items-start justify-between mb-3">
          <div>
            <h1
              className="text-xl font-serif font-bold mb-2"
              style={{ color: colors.emphasis }}
            >
              Interactive Flashcards
            </h1>
          </div>
          <button
            onClick={handleAddButton}
            className="flex items-center gap-1.5 px-2 py-1 rounded-lg font-medium transition-all hover:shadow-md flex-shrink-0"
            style={{
              background: 'linear-gradient(135deg, #9EC0FF 0%, #F8FAFF 45%, #FFF6E2 72%, #FFDF94 100%)',
              color: colors.emphasis,
            }}
          >
            <Plus size={20} />
            Add
          </button>
        </div>

        {/* Category Filters */}
        <div className="flex gap-2 flex-wrap overflow-x-auto pb-2">
          {categories.map(category => (
            <button
              key={category}
              onClick={() => handleCategoryChange(category)}
              className="px-4 py-2 rounded-full transition-all font-medium text-sm flex-shrink-0 whitespace-nowrap"
              style={{
                backgroundColor:
                  selectedCategory === category
                    ? colors.primary
                    : colors.surface,
                color:
                  selectedCategory === category
                    ? colors.text.light
                    : colors.text.secondary,
                border: `1px solid ${colors.border}`,
              }}
            >
              {category}
            </button>
          ))}
        </div>
      </div>

      {/* Content */}
      <div
        className="p-8 space-y-6"
        style={{ backgroundColor: colors.background }}
      >
        {/* Flashcard Grid */}
        <div className="grid grid-cols-3 gap-4">
          {pagedCards.map(card => (
            <FlashcardCard
              key={card.id}
              card={card}
              onCardClick={handleCardClick}
              onDelete={handleRequestDelete}
            />
          ))}
        </div>
        {filteredCards.length === 0 && (
          <div className="text-center py-12">
            <p style={{ color: colors.accent }}>No flashcards in this category yet.</p>
          </div>
        )}
        <Pagination currentPage={safePage} totalPages={totalPages} onPageChange={setCurrentPage} />
      </div>

      {/* Modal */}
      <FlashcardDetailModal
        card={selectedCard}
        isOpen={isModalOpen}
        onClose={handleCloseModal}
        onSave={handleSaveCard}
      />

      {/* Delete Confirmation */}
      <ConfirmDeleteModal
        card={pendingDeleteCard}
        isDeleting={isDeleting}
        onConfirm={handleConfirmDelete}
        onCancel={handleCancelDelete}
      />

      {/* Create Modal */}
      <CreateFlashcardModal
        isOpen={isCreateModalOpen}
        onClose={() => setIsCreateModalOpen(false)}
        onCreate={handleCreateCard}
      />
    </div>
  );
}
