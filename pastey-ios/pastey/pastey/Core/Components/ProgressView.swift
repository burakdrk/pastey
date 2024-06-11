//
//  ProgressView.swift
//  pastey
//
//  Created by Burak Duruk on 2024-06-10.
//

import SwiftUI

struct ProgressView: View {
    @State private var pulsate = false

    var body: some View {
        VStack {
            Text("pastey")
                .font(.largeTitle)
                .fontWeight(.bold)
                .padding(.bottom, 30)
                .scaleEffect(pulsate ? 1.4 : 1.2)
                .opacity(pulsate ? 0.5 : 1.0)
                .onAppear {
                    withAnimation(
                        Animation.easeInOut(duration: 1.0)
                            .repeatForever(autoreverses: true)
                    ) {
                        pulsate.toggle()
                    }
                }
        }
    }
}

#Preview {
    ProgressView()
}
