import React, { useRef } from 'react';
import { Calendar } from 'lucide-react';

interface DateTimePickerProps {
  value: string;
  onChange: (value: string) => void;
  language: string;
  placeholder?: string;
}

const DateTimePicker: React.FC<DateTimePickerProps> = ({ value, onChange, language, placeholder }) => {
  const inputRef = useRef<HTMLInputElement>(null);

  const handleClick = () => {
    inputRef.current?.showPicker();
  };

  return (
    <div 
      className="relative flex items-center bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg px-3 py-2 cursor-pointer hover:border-blue-500 dark:hover:border-blue-500 transition-colors group min-w-[180px]"
      onClick={handleClick}
    >
      <Calendar 
        size={16} 
        className="text-slate-400 group-hover:text-blue-500 transition-colors mr-2 shrink-0" 
      />
      <div className="flex-1 relative h-5 overflow-hidden">
        {!value && (
          <span className="absolute inset-0 text-slate-400 text-sm flex items-center pointer-events-none">
            {placeholder || 'Select date...'}
          </span>
        )}
        <input
          ref={inputRef}
          key={`date-${language}`}
          lang={language === 'zh' ? 'zh-CN' : 'en-US'}
          type="datetime-local"
          value={value}
          onChange={(e) => onChange(e.target.value)}
          className={`
            absolute inset-0 w-full h-full bg-transparent border-none outline-none text-sm 
            ${!value ? 'text-transparent' : 'text-slate-700 dark:text-slate-200'}
            [&::-webkit-calendar-picker-indicator]:hidden
            cursor-pointer
          `}
          style={{
            // Hide the default calendar icon in WebKit browsers
            WebkitAppearance: 'none',
          }}
        />
      </div>
    </div>
  );
};

export default DateTimePicker;
