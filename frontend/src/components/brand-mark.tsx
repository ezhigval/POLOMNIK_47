type BrandMarkProps = {
  className?: string;
};

export function BrandMark({ className = "size-9 bg-brand-800" }: BrandMarkProps) {
  return (
    <span
      className={`flex shrink-0 items-center justify-center rounded-full text-white ${className}`}
      aria-hidden
    >
      <svg viewBox="0 0 24 24" className="size-[55%]" fill="none" stroke="currentColor" strokeWidth="2.2">
        <path d="M12 4v16M8 8h8" strokeLinecap="round" />
      </svg>
    </span>
  );
}
