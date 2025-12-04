"use client"

import { HeroSection } from "@/app/(landing)/_sections/HeroSection"
import { FeaturesSection } from "@/app/(landing)/_sections/FeaturesSection"
import { CategoriesSection } from "@/app/(landing)/_sections/CategoriesSection"
import { CTASection } from "@/app/(landing)/_sections/CTASection"
import { PopularCoursesSection } from "@/app/(landing)/_sections/PopularCoursesSection"
import { TestimonialsSection } from "@/app/(landing)/_sections/TestimonialsSection"

export default function HomePage() {
  return (
    <>
      <HeroSection />
      <FeaturesSection />
      <CategoriesSection />
      <CTASection />
      <PopularCoursesSection />
      <TestimonialsSection />
    </>
  )
}
