export function HeroBlobs() {
  return (
    <div className="absolute inset-0 overflow-hidden pointer-events-none -z-10">
      <div className="absolute -top-20 -left-20 w-72 h-72 bg-primary/10 rounded-full blur-3xl animate-blob" />
      <div className="absolute -top-10 right-0 w-96 h-96 bg-emerald-200/20 rounded-full blur-3xl animate-blob-delay" />
      <div className="absolute top-40 left-1/3 w-64 h-64 bg-teal-200/15 rounded-full blur-3xl animate-blob-delay-2" />
    </div>
  );
}

export function AuthBlobs() {
  return (
    <div className="absolute inset-0 overflow-hidden pointer-events-none -z-10">
      <div className="absolute top-0 left-1/4 w-64 h-64 bg-primary/8 rounded-full blur-3xl animate-blob" />
      <div className="absolute bottom-0 right-1/4 w-80 h-80 bg-emerald-300/10 rounded-full blur-3xl animate-blob-delay" />
    </div>
  );
}

export function PageBlobs() {
  return (
    <div className="absolute inset-0 overflow-hidden pointer-events-none -z-10">
      <div className="absolute -top-32 -right-32 w-96 h-96 bg-primary/5 rounded-full blur-3xl" />
      <div className="absolute bottom-0 -left-20 w-72 h-72 bg-teal-200/8 rounded-full blur-3xl" />
    </div>
  );
}
