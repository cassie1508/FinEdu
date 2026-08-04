import { api } from '../../../lib/api';
import { LearningResource } from './types';

interface ApiResponse<T> {
  data: T;
}

const BASE_PATH = '/api/v1/learning_center/resources';

export async function fetchLearningResources(): Promise<LearningResource[]> {
  const res = await api.get<ApiResponse<LearningResource[]>>(BASE_PATH);
  return res.data;
}
