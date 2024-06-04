//
//  Color.swift
//  pastey
//
//  Created by Burak Duruk on 2024-06-03.
//

import Foundation
import SwiftUI

extension Color {
    static let theme = ColorTheme()
}

struct ColorTheme {
    let accent = Color("AccentColor")
    let background = Color("BackgroundColor")
    let secondaryText = Color("SecondaryTextColor")
}
