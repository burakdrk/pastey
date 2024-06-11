//
//  BigButton.swift
//  pastey
//
//  Created by Burak Duruk on 2024-06-10.
//

import SwiftUI

struct BigButton: View {
    let text: String
    let action: () -> Void
    var isFetching: Bool = false

    var body: some View {
        Button(action: action) {
            Text(text)
                .font(.headline)
                .foregroundColor(.white)
                .padding()
                .frame(maxWidth: .infinity)
                .background(!isFetching ? Color.theme.accent : Color.theme.accent.opacity(0.5))
                .cornerRadius(10)
        }
        .disabled(isFetching)
    }
}

#Preview {
    Group {
        BigButton(text: "Log in", action: {})
            .previewLayout(.sizeThatFits)

        BigButton(text: "Log in", action: {})
            .previewLayout(.sizeThatFits)
            .preferredColorScheme(.dark)
    }
}
