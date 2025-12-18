import React from 'react';
import Button from '~/presentation/components/common/Button';

export default function Home() {
  return (
    <div className="w-screen px-4 md:px-30">
      {/* Hero Section */}
      <div className="flex flex-col">
        <div className="flex w-full flex-col items-center justify-between md:flex-row">
          {/* HERO LEFT */}
          <div className="item-center flex w-full flex-col justify-center space-y-8 py-14 md:flex-1/2 md:py-0">
            <div className="item-center flex flex-col gap-12">
              <h1 className="text-text-primary text-5xl">
                From <span className="text-primary-3">Beginner to Pro__</span>
                <span className="text-primary-2">Master</span> Tech Skill with
                <span className="text-primary-3">Us</span>
              </h1>
              <p className="text-text-secondary">
                Online and On-Campus courses designed by industry experts. Start your journey in Web
                Development, Data Science, Cybersecurity, and more.
              </p>
            </div>

            <div className="flex flex-col gap-4 md:flex-row">
              <Button variant="fill" color="primary" size="lg">
                Explore Courses
              </Button>
              <Button variant="outline" color="secondary" size="lg">
                Get Started for Free
              </Button>
            </div>
          </div>

          {/* HERO RIGHT */}
          <div className="align-center flex w-100 justify-center md:flex-1/2">
            <img className="aspect-auto md:h-130" src="/illustrations/hero-right.png" alt="" />
          </div>
        </div>
        <div className="jusitfy-start flex flex-row py-6 md:py-0">
          <div className="ml-auto flex w-110 flex-col gap-6 md:flex-row">
            <div className="relative flex h-15 w-48 flex-row">
              <img className="absolute right-20" src="/peoples/Ellipse 1.png" alt="a person" />
              <img className="absolute right-10" src="/peoples/Ellipse 2.png" alt="a person" />
              <img className="absolute right-0" src="/peoples/Ellipse 3.png" alt="a person" />
            </div>
            <p className="text-text-secondary text-right text-xl font-medium">
              Learn alongside our community of 450,000+ students, instructors, and mentors
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
