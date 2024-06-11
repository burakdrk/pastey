//
//  HomeView.swift
//  pastey
//
//  Created by Burak Duruk on 2024-06-03.
//

import SwiftUI

struct HomeView: View {
    @EnvironmentObject var auth: AuthViewModel
    
    var body: some View {
        
        NavigationStack {
            VStack {
                HStack {
                    Text("pastey")
                        .font(.title)
                        .bold()
                    Spacer()
                    
                    Button(action: {
                        print("edit list")
                    }, label: {
                        Image(systemName: "square.and.pencil")
                            .font(.title2)
                    })
                }
                .padding(.horizontal, 20)
                .padding(.vertical, 10)
                

                List {
                    ForEach(0..<8) { _ in
                        EntryRowView(entry: DeveloperPreview.instance.entry)
                    }
                }.listStyle(PlainListStyle())
                
                HStack {
                    BigButton(text: "Copy", action: {
                        print("test")
                    })
                    
                    BigButton(text: "Paste", action: {
                        print("test")
                    })
                }.padding()
                
                
            }
            .navigationBarHidden(true)
        }
    }
}

#Preview {
    HomeView()
        .environmentObject(AuthViewModel())
}
