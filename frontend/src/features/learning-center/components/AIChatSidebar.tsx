import { useState } from 'react';
import { Send, Paperclip, MessageSquare, X } from 'lucide-react';
import { colors } from '../lib/colors';

interface ChatMessage {
  id: string;
  role: 'user' | 'assistant';
  content: string;
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
  const [messages, setMessages] = useState<ChatMessage[]>([
    {
      id: '1',
      role: 'assistant',
      content: "Hi, I'm your AI finance tutor. Ask me anything or upload a document, report, earnings report or SEC filing and I'll explain it.",
    },
  ]);
  const [input, setInput] = useState('');

  const handleSendMessage = () => {
    if (input.trim()) {
      const newMessage: ChatMessage = {
        id: Date.now().toString(),
        role: 'user',
        content: input,
      };
      setMessages([...messages, newMessage]);
      setInput('');

      // Simulate assistant response
      setTimeout(() => {
        setMessages(prev => [
          ...prev,
          {
            id: (Date.now() + 1).toString(),
            role: 'assistant',
            content: 'I understand your question. Let me break this down for you...',
          },
        ]);
      }, 500);
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
          AI-powered educational assistant
        </p>
      </div>

      {/* Conversation Area */}
      <div className="flex-1 overflow-y-auto p-6 space-y-4">
        {messages.map(msg => (
          <div
            key={msg.id}
            className={`flex ${msg.role === 'user' ? 'justify-end' : 'justify-start'}`}
          >
            <div
              className={`max-w-xs p-3 rounded-lg text-sm leading-relaxed ${
                msg.role === 'user'
                  ? 'rounded-br-none'
                  : 'rounded-bl-none'
              }`}
              style={{
                backgroundColor: msg.role === 'user' ? colors.primary : colors.secondary,
                color: msg.role === 'user' ? colors.text.light : colors.text.primary,
              }}
            >
              {msg.content}
            </div>
          </div>
        ))}

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
          <button
            className="p-1 transition-opacity hover:opacity-70"
            style={{ color: colors.accent }}
          >
            <Paperclip size={18} />
          </button>
          <input
            type="text"
            placeholder="Ask about any finance topic..."
            value={input}
            onChange={e => setInput(e.target.value)}
            onKeyDown={e => e.key === 'Enter' && handleSendMessage()}
            className="flex-1 bg-transparent outline-none text-sm"
            style={{ color: colors.text.primary }}
          />
          <button
            onClick={handleSendMessage}
            className="p-1 transition-opacity hover:opacity-70"
            style={{ color: colors.primary }}
          >
            <Send size={18} />
          </button>
        </div>
      </div>
    </div>
  );
}
