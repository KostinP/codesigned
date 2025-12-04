import { testimonials } from "@/lib/mocks/testimonials"
import { TestimonialCard } from "@/components/cards/TestimonialCard"

export function TestimonialsSection() {
  return (
    <section className="py-20 px-4 bg-background">
      <div className="max-w-7xl mx-auto flex flex-col gap-12 text-center">
        <h2 className="text-2xl md:text-3xl font-bold">Что говорят ученики</h2>

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 text-left">
          {testimonials.map((t, i) => (
            <TestimonialCard key={i} {...t} />
          ))}
        </div>
      </div>
    </section>
  )
}
