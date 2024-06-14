//
//  GradientHeaderView.swift
//  pastey
//
//  Created by Burak Duruk on 2024-06-13.
//

import SwiftUI
import FluidGradient

struct GradientHeaderView: View {
    var title: String
    
    var body: some View {
        ZStack(alignment: .top) {
            FluidGradient(blobs: [.theme.accent.opacity(0.5), .theme.accent, .theme.accent.opacity(1)], highlights: [.theme.accent.opacity(0.6), .theme.accent, .theme.accent.opacity(1)], speed: 0.4, blur: 0.8)
                .background(Color.theme.accent.opacity(0.3))
                .edgesIgnoringSafeArea(.top)
                .frame(height: 200)

            Text(title)
                .scaledToFit()
                .frame(height: 100)
                .padding()
                .shadow(radius: 15)
                .offset(y: 20)
                .font(.system(size: 50, weight: .bold))
                .foregroundColor(.white)
        }
    }
}

#Preview {
    GradientHeaderView(title: "Test")
}
