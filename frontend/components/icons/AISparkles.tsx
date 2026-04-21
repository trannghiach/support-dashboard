'use client';

import React from 'react';
import Icon from '@ant-design/icons';
import type { CustomIconComponentProps } from '@ant-design/icons/lib/components/Icon';

const RawSvg = (props: React.SVGProps<SVGSVGElement>) => (
    <svg
        viewBox="0 0 24 24"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
        {...props}
    >
        <style>
            {`
        @keyframes ai-pulse {
          0%, 100% { opacity: 1; transform: scale(1); }
          50% { opacity: 0.8; transform: scale(0.95); }
        }
        @keyframes ai-spin {
          from { transform: rotate(0deg); }
          to { transform: rotate(-360deg); }
        }
        .ai-wrapper {
          animation: ai-pulse 2.5s ease-in-out infinite;
          transform-origin: center;
        }
        .ai-star-spin {
          animation: ai-spin 8s linear infinite;
          transform-origin: 20px 6px;
        }
      `}
        </style>

        <g className="ai-wrapper">
            <defs>
                <linearGradient id="ai-magic-gradient" x1="0%" y1="0%" x2="100%" y2="100%">
                    <stop offset="0%" stopColor="#8B5CF6" />
                    <stop offset="50%" stopColor="#D946EF" />
                    <stop offset="100%" stopColor="#F97316" />
                </linearGradient>
            </defs>

            <path
                d="M10 2 C 10 6.5 13.5 10 18 10 C 13.5 10 10 13.5 10 18 C 10 13.5 6.5 10 2 10 C 6.5 10 10 6.5 10 2 Z"
                fill="url(#ai-magic-gradient)"
            />
            <path
                className="ai-star-spin"
                d="M20 2 C 20 4.5 21.5 6 24 6 C 21.5 6 20 7.5 20 10 C 20 7.5 18.5 6 16 6 C 18.5 6 20 4.5 20 2 Z"
                fill="url(#ai-magic-gradient)"
            />
            <path
                d="M18 16 C 18 17.5 19 18.5 20.5 18.5 C 19 18.5 18 19.5 18 21 C 18 19.5 17 18.5 15.5 18.5 C 17 18.5 18 17.5 18 16 Z"
                fill="url(#ai-magic-gradient)"
            />
        </g>
    </svg>
);

export const AISparklesIcon = (props: Partial<CustomIconComponentProps>) => (
    <Icon component={RawSvg} {...props} />
);