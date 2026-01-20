"use client";

import { useRouter } from "next/navigation";

export default function Home() {
  const router = useRouter();

  function handleStart() {
    router.push("/swipe");
  }

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-black p-4 relative overflow-hidden">
      {/* Animated background gradients */}
      <div className="absolute inset-0 bg-gradient-to-br from-purple-900/20 via-black to-cyan-900/20"></div>
      <div className="absolute top-1/4 left-1/4 w-96 h-96 bg-purple-500/10 rounded-full blur-3xl animate-pulse"></div>
      <div className="absolute bottom-1/4 right-1/4 w-96 h-96 bg-cyan-500/10 rounded-full blur-3xl animate-pulse delay-700"></div>

      <div className="text-center z-10 space-y-8">
        <div className="text-9xl animate-bounce drop-shadow-[0_0_25px_rgba(255,0,255,0.8)]">🍕</div>
        
        <h1 className="text-7xl font-black bg-gradient-to-r from-purple-400 via-pink-500 to-cyan-400 bg-clip-text text-transparent drop-shadow-[0_0_30px_rgba(168,85,247,0.8)] tracking-tight">
          BITE DECIDE
        </h1>
        
        <p className="text-gray-400 text-xl max-w-md mx-auto leading-relaxed">
          Swipe through food suggestions and let AI guide your cravings
        </p>
        
        <button
          onClick={handleStart}
          className="group relative px-12 py-5 bg-gradient-to-r from-purple-600 to-pink-600 rounded-full font-bold text-xl uppercase tracking-widest overflow-hidden transition-all duration-300 hover:scale-110 hover:shadow-[0_0_40px_rgba(168,85,247,0.8)]"
        >
          <span className="relative z-10">Let's Go!</span>
          <div className="absolute inset-0 bg-gradient-to-r from-pink-600 to-purple-600 opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
        </button>
      </div>
    </div>
  );
}
