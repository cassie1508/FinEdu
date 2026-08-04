import { ChevronLeft, ChevronRight } from 'lucide-react';
import { colors } from '../lib/colors';

interface PaginationProps {
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
}

const MAX_VISIBLE_PAGES = 6;

export function Pagination({ currentPage, totalPages, onPageChange }: PaginationProps) {
  if (totalPages <= 1) return null;

  const windowStart = Math.floor((currentPage - 1) / MAX_VISIBLE_PAGES) * MAX_VISIBLE_PAGES + 1;
  const windowEnd = Math.min(windowStart + MAX_VISIBLE_PAGES - 1, totalPages);
  const pages = Array.from({ length: windowEnd - windowStart + 1 }, (_, i) => windowStart + i);

  return (
    <div className="flex items-center justify-center gap-1.5 pt-2" aria-label="Pagination">
      <button
        onClick={() => onPageChange(currentPage - 1)}
        disabled={currentPage === 1}
        className="p-1.5 rounded-lg transition-all disabled:opacity-40 disabled:cursor-not-allowed hover:enabled:shadow-md"
        style={{ color: colors.text.secondary }}
        aria-label="Previous page"
      >
        <ChevronLeft size={16} />
      </button>

      {pages.map(page => (
        <button
          key={page}
          onClick={() => onPageChange(page)}
          aria-current={page === currentPage ? 'page' : undefined}
          className="min-w-8 h-8 px-2 rounded-lg text-sm font-medium transition-all"
          style={{
            backgroundColor: page === currentPage ? colors.primary : 'transparent',
            color: page === currentPage ? colors.text.light : colors.text.secondary,
          }}
        >
          {page}
        </button>
      ))}

      <button
        onClick={() => onPageChange(currentPage + 1)}
        disabled={currentPage === totalPages}
        className="p-1.5 rounded-lg transition-all disabled:opacity-40 disabled:cursor-not-allowed hover:enabled:shadow-md"
        style={{ color: colors.text.secondary }}
        aria-label="Next page"
      >
        <ChevronRight size={16} />
      </button>
    </div>
  );
}
