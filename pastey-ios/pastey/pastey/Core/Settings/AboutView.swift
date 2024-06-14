//
//  AboutView.swift
//  pastey
//
//  Created by Burak Duruk on 2024-06-13.
//

import SwiftUI

struct AboutView: View {
    let githubURL = URL(string: "https://github.com/burakdrk/pastey")!
    let defaultURL = URL(string: "https://burakduruk.com")!
    
    var body: some View {
        NavigationView {
            ZStack {
                List {
                    pasteySection
                        .listRowBackground(Color.theme.background.opacity(0.5))
                    applicationSection
                        .listRowBackground(Color.theme.background.opacity(0.5))
                }
            }
            .font(.headline)
            .listStyle(GroupedListStyle())
            .navigationTitle("About")
        }
    }
}

#Preview {
    AboutView()
}

extension AboutView {
    private var pasteySection: some View {
        Section(header: Text("pastey")) {
            VStack(alignment: .leading) {
                Image(uiImage: Bundle.main.icon ?? UIImage())
                    .resizable()
                    .frame(width: 100, height: 100)
                    .clipShape(RoundedRectangle(cornerRadius: 20))
                    .padding(.bottom, 15)
                Text("This app is an end-to-end encrypted clipboard manager. Built with SwiftUI.")
                    .font(.callout)
                    .fontWeight(.medium)
            }
            .padding(.vertical)
        }
    }
    
    private var applicationSection: some View {
        Section(header: Text("Application")) {
            Link("Source Code", destination: githubURL)
            Link("App Website", destination: defaultURL)
            Link("Developer", destination: defaultURL)
            Text("pastey \(Bundle.main.releaseVersionNumber ?? "") (\(Bundle.main.buildVersionNumber ?? ""))").foregroundColor(.theme.secondaryText)
        }
    }
}
