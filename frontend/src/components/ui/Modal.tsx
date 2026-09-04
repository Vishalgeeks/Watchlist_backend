import { ReactNode } from 'react';

interface ModalProps {
  isOpen: boolean;
  onClose: () => void;
  title?: string;
  children: ReactNode;
  maxWidth?: string;
}

export default function Modal({ isOpen, onClose, title, children, maxWidth = 'max-w-md' }: ModalProps) {
  if (!isOpen) return null;
  return (
    <div className="fixed inset-0 z-[90] flex items-center justify-center p-4">
      <div className="absolute inset-0 bg-black/60 backdrop-blur-sm" onClick={onClose} />
      <div className={`glass-panel relative w-full ${maxWidth} rounded-xl p-6 shadow-2xl`}>
        {title && (
          <div className="flex items-center justify-between mb-4">
            <h3 className="font-headline-md text-headline-md font-bold text-on-surface">{title}</h3>
            <button onClick={onClose} className="text-on-surface-variant hover:text-on-surface transition-colors">
              <span className="material-symbols-outlined">close</span>
            </button>
          </div>
        )}
        {children}
      </div>
    </div>
  );
}