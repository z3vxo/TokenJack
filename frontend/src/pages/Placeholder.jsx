export default function Placeholder({ title }) {
  return (
    <div className="m-auto max-w-md text-center">
      <h1 className="text-xl font-semibold tracking-[-0.02em]">{title}</h1>
      <p className="mt-2 text-sm leading-6 text-muted-foreground">
        This section is not built yet.
      </p>
    </div>
  );
}
