import React from 'react';

export default function LoadingState({ count = 5 }) {
  return (
    <div className="skeleton-list">
      {Array.from({ length: count }).map((_, idx) => (
        <div key={idx} className="skeleton-card">
          <div className="skeleton-cover pulse"></div>
          <div className="skeleton-details">
            <div className="skeleton-title pulse"></div>
            <div className="skeleton-author pulse"></div>
            <div className="skeleton-meta pulse"></div>
            <div className="skeleton-tags">
              <div className="skeleton-tag pulse"></div>
              <div className="skeleton-tag pulse"></div>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}
