"use client";

import { useEffect, useState, useRef } from "react";
import { useRouter } from "next/navigation";
import { createSession, getRecommendation, sendSwipe } from "../lib/api";

export default function SwipePage() {
  const router = useRouter();
  const [food, setFood] = useState<{
    name: string;
    description?: string;
  } | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [dragX, setDragX] = useState(0);
  const [isDragging, setIsDragging] = useState(false);
  const startX = useRef(0);

  const [likes, setLikes] = useState(0);
  const [dislikes, setDislikes] = useState(0);

  useEffect(() => {
    initSession();
  }, []);

  async function initSession() {
    try {
      setLoading(true);
      setError(null);

      let sessionId = localStorage.getItem("session_id");

      if (!sessionId) {
        const { session_id } = await createSession();
        sessionId = session_id;
        localStorage.setItem("session_id", sessionId);
      }

      try {
        const recommendation = await getRecommendation(sessionId);
        setFood(recommendation);
      } catch {
        localStorage.removeItem("session_id");
        const { session_id } = await createSession();
        localStorage.setItem("session_id", session_id);
        const recommendation = await getRecommendation(session_id);
        setFood(recommendation);
      }
    } catch (err) {
      setError("Something went wrong. Please try again.");
    } finally {
      setLoading(false);
    }
  }

  async function handleSwipe(action: "left" | "right" | "super") {
    const sessionId = localStorage.getItem("session_id");
    if (!sessionId || !food) return;

    if (action === "right" || action === "super") {
      setLikes(l => l + 1);
    } else if (action === "left") {
      setDislikes(d => d + 1);
    }

    try {
      setLoading(true);
      setError(null);

      await sendSwipe(sessionId, food.name, action);

      if (action === "super") {
        localStorage.setItem("selected_food", food.name);
        router.push("/complete");
      } else {
        try {
          const recommendation = await getRecommendation(sessionId);
          setFood(recommendation);
        } catch {
          localStorage.setItem("selected_food", food.name);
          router.push("/complete");
        }
      }
    } catch (err) {
      setError("Something went wrong. Please try again.");
    } finally {
      setLoading(false);
    }
  }

  function handleDragStart(clientX: number) {
    if (loading) return;
    setIsDragging(true);
    startX.current = clientX;
  }

  function handleDragMove(clientX: number) {
    if (!isDragging) return;
    setDragX(clientX - startX.current);
  }

  function handleDragEnd() {
    if (!isDragging) return;
    setIsDragging(false);

    const threshold = 100;
    if (dragX > threshold) {
      handleSwipe("right");
    } else if (dragX < -threshold) {
      handleSwipe("left");
    }
    setDragX(0);
  }

  function getCardStyle() {
    const rotate = dragX * 0.1;
    return {
      transform: `translateX(${dragX}px) rotate(${rotate}deg)`,
      transition: isDragging ? "none" : "transform 0.3s ease",
    };
  }

  function getSwipeIndicator() {
    if (dragX > 50) return "right";
    if (dragX < -50) return "left";
    return null;
  }

  if (loading && !food) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-black">
        <div className="text-purple-400 animate-pulse text-2xl font-bold">Loading...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen flex flex-col items-center justify-center bg-black p-4">
        <p className="text-red-500 mb-4 text-xl font-bold drop-shadow-[0_0_10px_rgba(239,68,68,0.8)]">{error}</p>
        <button
          onClick={initSession}
          className="px-8 py-4 bg-gradient-to-r from-purple-600 to-pink-600 rounded-full font-bold hover:scale-110 transition-transform shadow-[0_0_20px_rgba(168,85,247,0.6)]"
        >
          Retry
        </button>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-black p-4 select-none relative overflow-hidden">
      {/* Animated background */}
      <div className="absolute inset-0 bg-gradient-to-br from-purple-900/10 via-black to-cyan-900/10"></div>
      <div className="absolute top-0 left-1/4 w-96 h-96 bg-purple-500/5 rounded-full blur-3xl animate-pulse"></div>
      <div className="absolute bottom-0 right-1/4 w-96 h-96 bg-cyan-500/5 rounded-full blur-3xl animate-pulse delay-1000"></div>

      {/* Stats HUD */}
      <div className="w-full max-w-md flex justify-around mb-8 z-10">
        <div className="flex flex-col items-center space-y-2 bg-gradient-to-br from-green-500/10 to-green-600/5 backdrop-blur-sm rounded-2xl px-6 py-3 border border-green-500/30 shadow-[0_0_20px_rgba(34,197,94,0.3)]">
          <span className="text-green-400 text-4xl font-black drop-shadow-[0_0_10px_rgba(34,197,94,0.8)]">{likes}</span>
          <span className="text-xs text-green-300/70 uppercase tracking-widest font-bold">Likes</span>
        </div>
        
        <div className="flex flex-col items-center space-y-2 bg-gradient-to-br from-cyan-500/10 to-cyan-600/5 backdrop-blur-sm rounded-2xl px-6 py-3 border border-cyan-500/30 shadow-[0_0_20px_rgba(6,182,212,0.3)]">
          <span className="text-cyan-400 text-4xl font-black drop-shadow-[0_0_10px_rgba(6,182,212,0.8)]">{likes + dislikes}</span>
          <span className="text-xs text-cyan-300/70 uppercase tracking-widest font-bold">Total</span>
        </div>
        
        <div className="flex flex-col items-center space-y-2 bg-gradient-to-br from-red-500/10 to-red-600/5 backdrop-blur-sm rounded-2xl px-6 py-3 border border-red-500/30 shadow-[0_0_20px_rgba(239,68,68,0.3)]">
          <span className="text-red-400 text-4xl font-black drop-shadow-[0_0_10px_rgba(239,68,68,0.8)]">{dislikes}</span>
          <span className="text-xs text-red-300/70 uppercase tracking-widest font-bold">Nopes</span>
        </div>
      </div>

      {/* Food Card */}
      <div
        className={`relative w-full max-w-md aspect-[3/4] mb-12 cursor-grab active:cursor-grabbing z-10 ${
          dragX > 50 ? 'shadow-[0_0_50px_rgba(34,197,94,0.6)]' : 
          dragX < -50 ? 'shadow-[0_0_50px_rgba(239,68,68,0.6)]' : 
          'shadow-[0_0_50px_rgba(168,85,247,0.4)]'
        }`}
        style={getCardStyle()}
        onMouseDown={(e) => handleDragStart(e.clientX)}
        onMouseMove={(e) => handleDragMove(e.clientX)}
        onMouseUp={handleDragEnd}
        onMouseLeave={handleDragEnd}
        onTouchStart={(e) => handleDragStart(e.touches[0].clientX)}
        onTouchMove={(e) => handleDragMove(e.touches[0].clientX)}
        onTouchEnd={handleDragEnd}
      >
        <div className={`absolute inset-0 rounded-3xl bg-gradient-to-br from-gray-900 to-black border-4 ${
          dragX > 50 ? 'border-green-500' : 
          dragX < -50 ? 'border-red-500' : 
          'border-purple-500'
        } backdrop-blur-xl flex flex-col items-center justify-center p-8 transition-colors duration-200`}>
          
          {/* Swipe indicators */}
          {getSwipeIndicator() === "right" && (
            <div className="absolute top-8 left-8 px-6 py-3 bg-green-500/20 border-4 border-green-400 text-green-400 rounded-2xl font-black text-5xl rotate-[-20deg] shadow-[0_0_30px_rgba(34,197,94,0.8)] backdrop-blur-sm">
              LIKE
            </div>
          )}
          {getSwipeIndicator() === "left" && (
            <div className="absolute top-8 right-8 px-6 py-3 bg-red-500/20 border-4 border-red-400 text-red-400 rounded-2xl font-black text-5xl rotate-[20deg] shadow-[0_0_30px_rgba(239,68,68,0.8)] backdrop-blur-sm">
              NOPE
            </div>
          )}

          <div className="text-7xl mb-6 drop-shadow-[0_0_20px_rgba(255,255,255,0.5)] animate-pulse">🍽️</div>
          
          <h2 className="text-4xl font-black text-center bg-gradient-to-r from-purple-400 via-pink-400 to-cyan-400 bg-clip-text text-transparent mb-4 leading-tight drop-shadow-[0_0_20px_rgba(168,85,247,0.6)] px-4">
            {food?.name}
          </h2>
          
          {food?.description && (
            <p className="text-gray-300 text-center text-base leading-relaxed max-w-xs px-6 mb-4">
              {food.description}
            </p>
          )}

          <p className="absolute bottom-6 text-gray-500 text-xs uppercase tracking-widest font-bold">
            {isDragging ? "Release to Decide" : "Drag or Tap Below"}
          </p>
        </div>
      </div>

      {/* Action Buttons */}
      <div className="flex items-center gap-6 z-10">
        <button
          onClick={() => handleSwipe("left")}
          disabled={loading}
          className="group w-20 h-20 rounded-full bg-gradient-to-br from-red-500 to-red-600 text-white text-4xl flex items-center justify-center disabled:opacity-50 transition-all duration-300 hover:scale-125 hover:rotate-12 shadow-[0_0_30px_rgba(239,68,68,0.6)] hover:shadow-[0_0_50px_rgba(239,68,68,0.9)] border-2 border-red-400"
          aria-label="Nope"
        >
          <span className="group-hover:scale-110 transition-transform">✕</span>
        </button>

        <button
          onClick={() => handleSwipe("super")}
          disabled={loading}
          className="group w-24 h-24 rounded-full bg-gradient-to-br from-yellow-400 via-orange-500 to-pink-500 text-white text-4xl flex items-center justify-center disabled:opacity-50 transition-all duration-300 hover:scale-125 shadow-[0_0_40px_rgba(251,146,60,0.8)] hover:shadow-[0_0_60px_rgba(251,146,60,1)] border-4 border-yellow-300 animate-pulse"
          aria-label="This is it!"
        >
          <span className="group-hover:scale-110 group-hover:rotate-180 transition-all duration-300">⭐</span>
        </button>

        <button
          onClick={() => handleSwipe("right")}
          disabled={loading}
          className="group w-20 h-20 rounded-full bg-gradient-to-br from-green-500 to-green-600 text-white text-4xl flex items-center justify-center disabled:opacity-50 transition-all duration-300 hover:scale-125 hover:rotate-[-12deg] shadow-[0_0_30px_rgba(34,197,94,0.6)] hover:shadow-[0_0_50px_rgba(34,197,94,0.9)] border-2 border-green-400"
          aria-label="Like"
        >
          <span className="group-hover:scale-110 transition-transform">♥</span>
        </button>
      </div>
    </div>
  );
}
