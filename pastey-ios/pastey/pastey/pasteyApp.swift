//
//  pasteyApp.swift
//  pastey
//
//  Created by Burak Duruk on 2024-05-30.
//

import SwiftUI

@main
struct pasteyApp: App {
    @StateObject var auth = AuthViewModel()

    var body: some Scene {
        WindowGroup {
            RootView()
                .environmentObject(auth)
                .tint(.theme.accent)
        }
    }
}
