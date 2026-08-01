import { useState, useRef } from 'react';
import { Send, Paperclip, MessageSquare, X, Upload, AlertCircle, CheckCircle } from 'lucide-react';
import { colors } from '../lib/colors';
import { api } from '../../../lib/api';

interface ChatMessage {
  id: string;
  role: 'user' | 'assistant';
  content: string;
  sources?: string[]; // filenames of source documents
}

interface UploadedDocument {
  id: string;
  filename: string;
  chunksCount: number;
  status: 'uploading' | 'success' | 'error';
  error?: string;
}

interface DocumentUploadResponse {
  success: boolean;
  documentId: string;
  chunksCount: number;
  embeddingsStored: number;
  embeddingModel: string;
  message: string;
}

interface RetrievedChunk {
  text: string;
  filename: string;
  similarityScore: number;
}

interface RAGQueryResponse {
  success: boolean;
  query: string;
  retrievedChunks: RetrievedChunk[];
  generatedResponse: string;
}

const suggestedPrompts = [
  'What is an ETF?',
  'Explain compound interest simply',
  'How much should I save for emergencies?',
];

interface AIChatSidebarProps {
  mode?: 'sidebar' | 'popup';
  onClose?: () => void;
}

export function AIChatSidebar({ mode = 'sidebar', onClose }: AIChatSidebarProps) {
  const isPopup = mode === 'popup';
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [sessionId] = useState(() => `session-${Date.now()}`);
  const [messages, setMessages] = useState<ChatMessage[]>([
    {
      id: '1',
      role: 'assistant',
      content: "Hi, I'm your AI finance tutor. Upload a PDF document (earnings report, SEC filing, etc.) and ask me questions about it. I'll explain it using AI!",
    },
  ]);
  const [input, setInput] = useState('');
  const [uploadedDocs, setUploadedDocs] = useState<UploadedDocument[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  const handleFileSelect = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    if (!file.name.toLowerCase().endsWith('.pdf')) {
      alert('Please upload a PDF file');
      return;
    }

    // Create document entry with uploading status
    const docId = Date.now().toString();
    const newDoc: UploadedDocument = {
      id: docId,
      filename: file.name,
      chunksCount: 0,
      status: 'uploading',
    };
    setUploadedDocs(prev => [...prev, newDoc]);

    try {
      // Create FormData for multipart upload
      const formData = new FormData();
      formData.append('sessionId', sessionId);
      formData.append('file', file);

      // Call document upload API with FormData
      const response = await fetch(`${import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080'}/api/v1/documents/upload`, {
        method: 'POST',
        body: formData,
      });

      if (!response.ok) {
        throw new Error(`Upload failed: ${response.status}`);
      }

      const result = await response.json() as DocumentUploadResponse;

      if (result.success) {
        setUploadedDocs(prev =>
          prev.map(doc =>
            doc.id === docId
              ? {
                  ...doc,
                  status: 'success',
                  chunksCount: result.chunksCount || 0,
                }
              : doc
          )
        );

        // Add system message about successful upload
        setMessages(prev => [
          ...prev,
          {
            id: `sys-${Date.now()}`,
            role: 'assistant',
            content: `✅ Successfully uploaded "${file.name}" (${result.chunksCount} chunks). Now you can ask questions about it!`,
          },
        ]);
      } else {
        throw new Error(result.message || 'Upload failed');
      }
    } catch (error: any) {
      setUploadedDocs(prev =>
        prev.map(doc =>
          doc.id === docId
            ? {
                ...doc,
                status: 'error',
                error: error.message || 'Upload failed',
              }
            : doc
        )
      );
    }

    // Clear file input
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  const handleSendMessage = async () => {
    if (input.trim() && !isLoading) {
      const userMessage: ChatMessage = {
        id: Date.now().toString(),
        role: 'user',
        content: input,
      };
      setMessages(prev => [...prev, userMessage]);
      setInput('');
      setIsLoading(true);

      try {
        // If documents are uploaded, query RAG pipeline
        if (uploadedDocs.some(doc => doc.status === 'success')) {
          const response = await api.post<RAGQueryResponse>('/api/v1/rag/query', {
            sessionId,
            query: userMessage.content,
            topK: 5,
          });

          if (response.success) {
            const sourceFiles = Array.from(
              new Set(
                (response.retrievedChunks || []).map((chunk: any) => chunk.filename).filter(Boolean)
              )
            );

            const assistantMessage: ChatMessage = {
              id: `${Date.now()}-assistant`,
              role: 'assistant',
              content: response.generatedResponse || 'No response generated',
              sources: sourceFiles.length > 0 ? sourceFiles : undefined,
            };
            setMessages(prev => [...prev, assistantMessage]);
          } else {
            throw new Error(response.generatedResponse || 'Query failed');
          }
        } else {
          // Fallback: generic response when no documents uploaded
          const assistantMessage: ChatMessage = {
            id: `${Date.now()}-assistant`,
            role: 'assistant',
            content:
              'I need a document to answer your question. Please upload a PDF first (earnings report, SEC filing, etc.), then ask your question.',
          };
          setMessages(prev => [...prev, assistantMessage]);
        }
      } catch (error: any) {
        const errorMessage: ChatMessage = {
          id: `${Date.now()}-error`,
          role: 'assistant',
          content: `Sorry, I encountered an error: ${error.message}`,
        };
        setMessages(prev => [...prev, errorMessage]);
      } finally {
        setIsLoading(false);
      }
    }
  };

  const handleSuggestedPrompt = (prompt: string) => {
    setInput(prompt);
  };

  return (
    <div
      className={`flex flex-col ${isPopup ? 'h-[70vh] max-h-[640px] rounded-2xl border shadow-2xl' : 'h-screen border-r'}`}
      style={{
        backgroundColor: colors.background,
        borderColor: colors.border,
        width: isPopup ? 'min(420px, calc(100vw - 32px))' : '320px',
      }}
    >
      {/* Header */}
      <div className="p-6 border-b" style={{ borderColor: colors.border }}>
        <div className="flex items-center justify-between mb-2">
          <div className="flex items-center gap-2">
            <MessageSquare size={24} style={{ color: colors.primary }} />
            <h1 className="text-lg font-semibold" style={{ color: colors.emphasis }}>
              AI Finance Tutor
            </h1>
          </div>
          {isPopup && (
            <button
              onClick={onClose}
              className="p-1 rounded-md transition-opacity hover:opacity-70"
              style={{ color: colors.accent }}
              aria-label="Close AI chat"
            >
              <X size={18} />
            </button>
          )}
        </div>
        <p className="text-sm" style={{ color: colors.accent }}>
          Upload PDFs & ask questions
        </p>

        {/* Uploaded Documents */}
        {uploadedDocs.length > 0 && (
          <div className="mt-4 space-y-2">
            {uploadedDocs.map(doc => (
              <div
                key={doc.id}
                className="flex items-center gap-2 text-xs p-2 rounded-lg"
                style={{ backgroundColor: colors.surface }}
              >
                {doc.status === 'uploading' && (
                  <div className="animate-spin" style={{ color: colors.primary }}>
                    <Upload size={14} />
                  </div>
                )}
                {doc.status === 'success' && (
                  <CheckCircle size={14} style={{ color: '#10b981' }} />
                )}
                {doc.status === 'error' && (
                  <AlertCircle size={14} style={{ color: '#ef4444' }} />
                )}
                <div className="flex-1 min-w-0">
                  <p className="truncate font-medium" style={{ color: colors.text.primary }}>
                    {doc.filename}
                  </p>
                  {doc.status === 'success' && (
                    <p className="text-xs" style={{ color: colors.accent }}>
                      {doc.chunksCount} chunks
                    </p>
                  )}
                  {doc.status === 'error' && (
                    <p className="text-xs" style={{ color: '#ef4444' }}>
                      {doc.error}
                    </p>
                  )}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Conversation Area */}
      <div className="flex-1 overflow-y-auto p-6 space-y-4">
        {messages.map(msg => (
          <div key={msg.id} className={`flex ${msg.role === 'user' ? 'justify-end' : 'justify-start'}`}>
            <div className="w-full max-w-xs">
              <div
                className={`p-3 rounded-lg text-sm leading-relaxed ${
                  msg.role === 'user' ? 'rounded-br-none' : 'rounded-bl-none'
                }`}
                style={{
                  backgroundColor: msg.role === 'user' ? colors.primary : colors.secondary,
                  color: msg.role === 'user' ? colors.text.light : colors.text.primary,
                }}
              >
                {msg.content}
              </div>
              {msg.sources && msg.sources.length > 0 && (
                <p className="text-xs mt-1 px-2" style={{ color: colors.accent }}>
                  📄 Source: {msg.sources.join(', ')}
                </p>
              )}
            </div>
          </div>
        ))}

        {/* Loading indicator */}
        {isLoading && (
          <div className="flex justify-start">
            <div
              className="p-3 rounded-lg rounded-bl-none"
              style={{ backgroundColor: colors.secondary }}
            >
              <div className="flex gap-1">
                <div
                  className="w-2 h-2 rounded-full animate-bounce"
                  style={{ backgroundColor: colors.primary, animationDelay: '0ms' }}
                />
                <div
                  className="w-2 h-2 rounded-full animate-bounce"
                  style={{ backgroundColor: colors.primary, animationDelay: '150ms' }}
                />
                <div
                  className="w-2 h-2 rounded-full animate-bounce"
                  style={{ backgroundColor: colors.primary, animationDelay: '300ms' }}
                />
              </div>
            </div>
          </div>
        )}

        {/* Suggested Prompts (show when limited messages) */}
        {messages.length <= 1 && (
          <div className="pt-4 space-y-2">
            <p className="text-xs font-medium" style={{ color: colors.accent }}>
              Suggested questions:
            </p>
            {suggestedPrompts.map(prompt => (
              <button
                key={prompt}
                onClick={() => handleSuggestedPrompt(prompt)}
                className="w-full text-left text-sm p-2 rounded-lg transition-colors hover:opacity-80"
                style={{
                  backgroundColor: colors.surface,
                  color: colors.text.primary,
                  border: `1px solid ${colors.border}`,
                }}
              >
                {prompt}
              </button>
            ))}
          </div>
        )}
      </div>

      {/* Input Area */}
      <div className="p-4 border-t" style={{ borderColor: colors.border }}>
        <div
          className="flex items-center gap-2 px-4 py-3 rounded-full border"
          style={{
            backgroundColor: colors.surface,
            borderColor: colors.border,
          }}
        >
          <input
            ref={fileInputRef}
            type="file"
            accept=".pdf"
            onChange={handleFileSelect}
            className="hidden"
            aria-label="Upload PDF"
          />
          <button
            onClick={() => fileInputRef.current?.click()}
            className="p-1 transition-opacity hover:opacity-70 flex-shrink-0"
            style={{ color: colors.accent }}
            title="Upload PDF document"
          >
            <Paperclip size={18} />
          </button>
          <input
            type="text"
            placeholder="Ask about the document..."
            value={input}
            onChange={e => setInput(e.target.value)}
            onKeyDown={e => e.key === 'Enter' && !isLoading && handleSendMessage()}
            disabled={isLoading}
            className="flex-1 bg-transparent outline-none text-sm disabled:opacity-50"
            style={{ color: colors.text.primary }}
          />
          <button
            onClick={handleSendMessage}
            disabled={isLoading}
            className="p-1 transition-opacity hover:opacity-70 disabled:opacity-50 flex-shrink-0"
            style={{ color: colors.primary }}
          >
            <Send size={18} />
          </button>
        </div>
      </div>
    </div>
  );
}
