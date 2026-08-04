import { api } from '../../../lib/api';
import { Flashcard } from './types';

interface ApiResponse<T> {
  data: T;
}

export interface FlashcardInput {
  title: string;
  category: string;
  whyItMatters: string;
  definition: string;
  example: string;
  commonMisconceptions: string;
}

const BASE_PATH = '/api/v1/learning/flashcards';

export async function fetchFlashcards(): Promise<Flashcard[]> {
  const res = await api.get<ApiResponse<Flashcard[]>>(BASE_PATH);
  return res.data ?? [];
}

export async function createFlashcard(input: FlashcardInput): Promise<Flashcard> {
  const res = await api.post<ApiResponse<Flashcard>>(BASE_PATH, input);
  return res.data;
}

export async function updateFlashcard(id: string, input: FlashcardInput): Promise<Flashcard> {
  const res = await api.put<ApiResponse<Flashcard>>(`${BASE_PATH}/${id}`, input);
  return res.data;
}

export async function deleteFlashcard(id: string): Promise<void> {
  await api.delete(`${BASE_PATH}/${id}`);
}

export async function reviewFlashcard(id: string): Promise<Flashcard> {
  const res = await api.post<ApiResponse<Flashcard>>(`${BASE_PATH}/${id}/review`, {});
  return res.data;
}
