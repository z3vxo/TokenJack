export default function BrandMark({ size = 36 }) {
  return (
    <span
      className="inline-flex items-center justify-center rounded-mark bg-primary text-on-primary"
      style={{ width: size, height: size }}
      aria-hidden="true"
    >
      <svg
        width={size * 0.5}
        height={size * 0.5}
        viewBox="0 0 24 24"
        fill="none"
        focusable="false"
      >
        <circle
          cx="12"
          cy="9.5"
          r="3.4"
          stroke="currentColor"
          strokeWidth="2.1"
        />
        <path
          d="M12 12.9V21"
          stroke="currentColor"
          strokeWidth="2.1"
          strokeLinecap="round"
        />
        <path
          d="M12 17.6h3.1"
          stroke="currentColor"
          strokeWidth="2.1"
          strokeLinecap="round"
        />
      </svg>
    </span>
  );
}
