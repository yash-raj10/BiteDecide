"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { getFoodInfo } from "../lib/api";

interface FoodInfo {
  description: string;
  nutrients: string;
  similar_foods: string[];
}

export default function CompletePage() {
  const router = useRouter();
  const [selectedFood, setSelectedFood] = useState<string | null>(null);
  const [foodInfo, setFoodInfo] = useState<FoodInfo | null>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    const food = localStorage.getItem("selected_food");
    setSelectedFood(food);

    if (food) {
      setLoading(true);
      getFoodInfo(food)
        .then(setFoodInfo)
        .catch(() => {})
        .finally(() => setLoading(false));
    }
  }, []);

  function handleNewSession() {
    localStorage.removeItem("session_id");
    localStorage.removeItem("selected_food");
    router.push("/");
  }

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-black p-4 relative overflow-hidden">
      {/* Animated background */}
      <div className="absolute inset-0 bg-gradient-to-br from-purple-900/20 via-black to-pink-900/20"></div>
      <div className="absolute top-1/3 left-1/3 w-96 h-96 bg-purple-500/10 rounded-full blur-3xl animate-pulse"></div>
      <div className="absolute bottom-1/3 right-1/3 w-96 h-96 bg-pink-500/10 rounded-full blur-3xl animate-pulse delay-700"></div>

      <div className="relative z-10 text-center space-y-8 max-w-2xl w-full">
        {/* Success animation */}
        <div className="text-9xl animate-bounce drop-shadow-[0_0_30px_rgba(251,146,60,0.8)]">🎉</div>
        
        <div className="space-y-4">
          <h1 className="text-4xl md:text-6xl font-black bg-gradient-to-r from-yellow-400 via-orange-500 to-pink-500 bg-clip-text text-transparent drop-shadow-[0_0_30px_rgba(251,146,60,0.8)]">
            PERFECT!
          </h1>
          
          <p className="text-lg md:text-xl text-gray-400">
            You picked
          </p>
          
          <div className="bg-gradient-to-br from-purple-500/20 to-pink-500/20 backdrop-blur-xl rounded-3xl p-6 md:p-8 border-2 border-purple-500/50 shadow-[0_0_40px_rgba(168,85,247,0.4)]">
            <p className="text-3xl md:text-5xl font-black bg-gradient-to-r from-purple-400 via-pink-400 to-cyan-400 bg-clip-text text-transparent leading-tight">
              {selectedFood}
            </p>
          </div>
          
          <p className="text-base md:text-lg text-gray-500 italic">
            Enjoy your meal! 🍽️
          </p>
        </div>

        {/* Food Info Section */}
        {loading && (
          <div className="text-cyan-400 animate-pulse text-lg">Loading info...</div>
        )}

        {foodInfo && !loading && (
          <div className="space-y-6 text-left bg-gradient-to-br from-gray-900/50 to-black/50 backdrop-blur-xl rounded-3xl p-6 md:p-8 border border-gray-700/50">
            {/* AI Generated Tag */}
            <div className="flex justify-center mb-2">
              <span className="px-4 py-1.5 bg-gradient-to-r from-purple-500/20 to-pink-500/20 border border-purple-400/30 text-purple-300 rounded-full text-xs font-bold uppercase tracking-wider">
                ✨ AI Generated
              </span>
            </div>

            <div>
              <h3 className="text-xl md:text-2xl font-bold text-purple-400 mb-3 flex items-center gap-2">
                <span>📖</span> About
              </h3>
              <p className="text-gray-300 text-sm md:text-base leading-relaxed">{foodInfo.description}</p>
            </div>

            <div>
              <h3 className="text-xl md:text-2xl font-bold text-green-400 mb-3 flex items-center gap-2">
                <span>🥗</span> Nutrients
              </h3>
              <p className="text-gray-300 text-sm md:text-base leading-relaxed">{foodInfo.nutrients}</p>
            </div>

            {foodInfo.similar_foods.length > 0 && (
              <div>
                <h3 className="text-xl md:text-2xl font-bold text-cyan-400 mb-3 flex items-center gap-2">
                  <span>🍽️</span> Similar Foods
                </h3>
                <div className="flex flex-wrap gap-2 md:gap-3">
                  {foodInfo.similar_foods.map((food) => (
                    <span
                      key={food}
                      className="px-3 md:px-4 py-1.5 md:py-2 bg-gradient-to-r from-cyan-500/20 to-blue-500/20 border border-cyan-500/30 text-cyan-300 rounded-full text-xs md:text-sm font-medium"
                    >
                      {food}
                    </span>
                  ))}
                </div>
              </div>
            )}
          </div>
        )}

        <button
          onClick={handleNewSession}
          className="group w-full px-10 py-5 bg-gradient-to-r from-purple-600 to-pink-600 rounded-full font-bold text-xl uppercase tracking-widest overflow-hidden transition-all duration-300 hover:scale-110 hover:shadow-[0_0_40px_rgba(168,85,247,0.8)] mt-8"
        >
          <span className="relative z-10">Start New Session</span>
          <div className="absolute inset-0 bg-gradient-to-r from-pink-600 to-purple-600 opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
        </button>
      </div>
    </div>
  );
}
